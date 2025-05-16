package create

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// generateReactProject creates a new React project
func generateReactProject(stateManagement string, options map[string]string) (string, error) {
	// Get project name from options or use a default
	projectName := options["name"]
	if projectName == "" {
		projectName = "my-react-app"
	}

	// Check if Node.js is installed
	if err := checkNodeInstalled(); err != nil {
		return "", err
	}

	// Create the project using create-react-app
	if err := createBaseReactProject(projectName); err != nil {
		return "", err
	}

	// Set up the project structure based on state management
	switch strings.ToLower(stateManagement) {
	case "redux":
		if err := setupReactReduxArchitecture(projectName); err != nil {
			return "", err
		}
	case "context":
		if err := setupReactContextAPIArchitecture(projectName); err != nil {
			return "", err
		}
	case "mobx":
		if err := setupReactMobXArchitecture(projectName); err != nil {
			return "", err
		}
	case "recoil":
		if err := setupReactRecoilArchitecture(projectName); err != nil {
			return "", err
		}
	default:
		// Default to a basic structure without specific state management
		if err := setupBasicReactArchitecture(projectName); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("✅ React project '%s' created successfully with %s architecture!",
		projectName,
		getReactArchitectureName(stateManagement)), nil
}

// createBaseReactProject creates a new React project using create-react-app
func createBaseReactProject(name string) error {
	// Use npx to run create-react-app without installing it globally
	cmd := exec.Command("npx", "create-react-app", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// setupBasicReactArchitecture sets up a basic React project structure
func setupBasicReactArchitecture(projectPath string) error {
	// Create additional directories for a clean architecture
	dirs := []string{
		"src/components",
		"src/hooks",
		"src/utils",
		"src/assets",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	// Create a sample utility function
	utilsPath := filepath.Join(projectPath, "src/utils", "helpers.js")
	utilsContent := `/**
 * Format a date string
 * @param {string} dateString - The date string to format
 * @returns {string} Formatted date string
 */
export function formatDate(dateString) {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(date);
}

/**
 * Truncate text to a specific length
 * @param {string} text - The text to truncate
 * @param {number} length - Maximum length
 * @returns {string} Truncated text
 */
export function truncateText(text, length = 100) {
  if (text.length <= length) return text;
  return text.slice(0, length) + '...';
}
`
	if err := os.WriteFile(utilsPath, []byte(utilsContent), 0644); err != nil {
		return fmt.Errorf("failed to create helpers.js: %w", err)
	}

	// Create a sample component
	buttonComponentPath := filepath.Join(projectPath, "src/components", "Button.jsx")
	buttonComponentContent := `import React from 'react';
import './Button.css';

/**
 * Button component with variants
 */
function Button({ children, variant = 'primary', onClick }) {
  const getButtonClass = () => {
    const baseClass = 'button';
    return variant ? baseClass + ' ' + baseClass + '--' + variant : baseClass;
  };

  return (
    <button className={getButtonClass()} onClick={onClick}>
      {children}
    </button>
  );
}

export default Button;
`
	if err := os.WriteFile(buttonComponentPath, []byte(buttonComponentContent), 0644); err != nil {
		return fmt.Errorf("failed to create Button.jsx: %w", err)
	}

	// Create CSS for the button component
	buttonCSSPath := filepath.Join(projectPath, "src/components", "Button.css")
	buttonCSSContent := `.button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s, opacity 0.3s;
}

.button:hover {
  opacity: 0.9;
}

.button--primary {
  background-color: #0070f3;
  color: white;
}

.button--secondary {
  background-color: #f3f3f3;
  color: #333;
}

.button--danger {
  background-color: #ff0000;
  color: white;
}
`
	if err := os.WriteFile(buttonCSSPath, []byte(buttonCSSContent), 0644); err != nil {
		return fmt.Errorf("failed to create Button.css: %w", err)
	}

	// Create a custom hook
	hookPath := filepath.Join(projectPath, "src/hooks", "useLocalStorage.js")
	hookContent := `import { useState, useEffect } from 'react';

/**
 * Custom hook for using localStorage with React state
 * @param {string} key - The localStorage key
 * @param {any} initialValue - The initial value
 * @returns {Array} [storedValue, setValue]
 */
function useLocalStorage(key, initialValue) {
  // Get from local storage then parse stored json or return initialValue
  const readValue = () => {
    if (typeof window === 'undefined') {
      return initialValue;
    }

    try {
      const item = window.localStorage.getItem(key);
      return item ? JSON.parse(item) : initialValue;
    } catch (error) {
      console.warn("Error reading localStorage key '" + key + "':", error);
      return initialValue;
    }
  };

  // State to store our value
  const [storedValue, setStoredValue] = useState(readValue);

  // Return a wrapped version of useState's setter function that persists the new value to localStorage
  const setValue = (value) => {
    try {
      // Allow value to be a function so we have same API as useState
      const valueToStore = value instanceof Function ? value(storedValue) : value;

      // Save state
      setStoredValue(valueToStore);

      // Save to local storage
      if (typeof window !== 'undefined') {
        window.localStorage.setItem(key, JSON.stringify(valueToStore));
      }
    } catch (error) {
      console.warn("Error setting localStorage key '" + key + "':", error);
    }
  };

  useEffect(() => {
    setStoredValue(readValue());
  }, []);

  return [storedValue, setValue];
}

export default useLocalStorage;
`
	if err := os.WriteFile(hookPath, []byte(hookContent), 0644); err != nil {
		return fmt.Errorf("failed to create useLocalStorage.js: %w", err)
	}

	return nil
}

// setupReactReduxArchitecture sets up a React project with Redux
func setupReactReduxArchitecture(projectPath string) error {
	// Install Redux dependencies
	cmd := exec.Command("npm", "install", "redux", "react-redux", "@reduxjs/toolkit")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Redux dependencies: %w", err)
	}

	// Create directories for Redux architecture
	dirs := []string{
		"src/components",
		"src/hooks",
		"src/utils",
		"src/assets",
		"src/store",
		"src/store/slices",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	// Create Redux store
	storePath := filepath.Join(projectPath, "src/store", "index.js")
	storeContent := `import { configureStore } from '@reduxjs/toolkit';
import counterReducer from './slices/counterSlice';

export const store = configureStore({
  reducer: {
    counter: counterReducer,
    // Add more reducers here
  },
});
`
	if err := os.WriteFile(storePath, []byte(storeContent), 0644); err != nil {
		return fmt.Errorf("failed to create store/index.js: %w", err)
	}

	// Create a sample Redux slice
	slicePath := filepath.Join(projectPath, "src/store/slices", "counterSlice.js")
	sliceContent := `import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  value: 0,
};

export const counterSlice = createSlice({
  name: 'counter',
  initialState,
  reducers: {
    increment: (state) => {
      state.value += 1;
    },
    decrement: (state) => {
      state.value -= 1;
    },
    incrementByAmount: (state, action) => {
      state.value += action.payload;
    },
  },
});

export const { increment, decrement, incrementByAmount } = counterSlice.actions;

export default counterSlice.reducer;
`
	if err := os.WriteFile(slicePath, []byte(sliceContent), 0644); err != nil {
		return fmt.Errorf("failed to create counterSlice.js: %w", err)
	}

	// Update index.js to include Redux provider
	indexPath := filepath.Join(projectPath, "src", "index.js")

	// Check if index.js exists
	_, err := os.Stat(indexPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to check index.js: %w", err)
	}

	// Create new content with Redux provider
	newIndexContent := `import React from 'react';
import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import { store } from './store';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <Provider store={store}>
      <App />
    </Provider>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
`

	// Write the updated content
	if err := os.WriteFile(indexPath, []byte(newIndexContent), 0644); err != nil {
		return fmt.Errorf("failed to update index.js: %w", err)
	}

	// Create a sample counter component
	counterComponentPath := filepath.Join(projectPath, "src/components", "Counter.jsx")
	counterComponentContent := `import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { increment, decrement, incrementByAmount } from '../store/slices/counterSlice';
import './Counter.css';

function Counter() {
  const count = useSelector((state) => state.counter.value);
  const dispatch = useDispatch();

  return (
    <div className="counter">
      <h2>Redux Counter</h2>
      <div className="counter-value">{count}</div>
      <div className="counter-buttons">
        <button onClick={() => dispatch(decrement())}>-</button>
        <button onClick={() => dispatch(increment())}>+</button>
        <button onClick={() => dispatch(incrementByAmount(5))}>+5</button>
      </div>
    </div>
  );
}

export default Counter;
`
	if err := os.WriteFile(counterComponentPath, []byte(counterComponentContent), 0644); err != nil {
		return fmt.Errorf("failed to create Counter.jsx: %w", err)
	}

	// Create CSS for the counter component
	counterCSSPath := filepath.Join(projectPath, "src/components", "Counter.css")
	counterCSSContent := `.counter {
  text-align: center;
  margin: 2rem auto;
  padding: 1rem;
  max-width: 300px;
  border: 1px solid #ccc;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.counter-value {
  font-size: 3rem;
  font-weight: bold;
  margin: 1rem 0;
}

.counter-buttons {
  display: flex;
  justify-content: center;
  gap: 0.5rem;
}

.counter-buttons button {
  padding: 0.5rem 1rem;
  font-size: 1.25rem;
  border: none;
  border-radius: 4px;
  background-color: #0070f3;
  color: white;
  cursor: pointer;
}

.counter-buttons button:hover {
  background-color: #0060df;
}
`
	if err := os.WriteFile(counterCSSPath, []byte(counterCSSContent), 0644); err != nil {
		return fmt.Errorf("failed to create Counter.css: %w", err)
	}

	return nil
}

// setupReactContextAPIArchitecture sets up a React project with Context API
func setupReactContextAPIArchitecture(projectPath string) error {
	// Create directories for Context API architecture
	dirs := []string{
		"src/components",
		"src/hooks",
		"src/utils",
		"src/assets",
		"src/contexts",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	// Create a sample context
	contextPath := filepath.Join(projectPath, "src/contexts", "CounterContext.jsx")
	contextContent := `import React, { createContext, useContext, useState } from 'react';

// Create the context
const CounterContext = createContext();

// Create a provider component
export function CounterProvider({ children }) {
  const [count, setCount] = useState(0);

  const increment = () => setCount(count + 1);
  const decrement = () => setCount(count - 1);
  const reset = () => setCount(0);
  const incrementByAmount = (amount) => setCount(count + amount);

  const value = {
    count,
    increment,
    decrement,
    reset,
    incrementByAmount,
  };

  return (
    <CounterContext.Provider value={value}>
      {children}
    </CounterContext.Provider>
  );
}

// Create a custom hook for using the context
export function useCounter() {
  const context = useContext(CounterContext);
  if (context === undefined) {
    throw new Error('useCounter must be used within a CounterProvider');
  }
  return context;
}
`
	if err := os.WriteFile(contextPath, []byte(contextContent), 0644); err != nil {
		return fmt.Errorf("failed to create CounterContext.jsx: %w", err)
	}

	// Update index.js to include Context provider
	indexPath := filepath.Join(projectPath, "src", "index.js")

	// Check if index.js exists
	_, err := os.Stat(indexPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to check index.js: %w", err)
	}

	// Create new content with Context provider
	newIndexContent := `import React from 'react';
import ReactDOM from 'react-dom/client';
import { CounterProvider } from './contexts/CounterContext';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <CounterProvider>
      <App />
    </CounterProvider>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
`

	// Write the updated content
	if err := os.WriteFile(indexPath, []byte(newIndexContent), 0644); err != nil {
		return fmt.Errorf("failed to update index.js: %w", err)
	}

	// Create a sample counter component using Context
	counterComponentPath := filepath.Join(projectPath, "src/components", "Counter.jsx")
	counterComponentContent := `import React from 'react';
import { useCounter } from '../contexts/CounterContext';
import './Counter.css';

function Counter() {
  const { count, increment, decrement, incrementByAmount } = useCounter();

  return (
    <div className="counter">
      <h2>Context API Counter</h2>
      <div className="counter-value">{count}</div>
      <div className="counter-buttons">
        <button onClick={decrement}>-</button>
        <button onClick={increment}>+</button>
        <button onClick={() => incrementByAmount(5)}>+5</button>
      </div>
    </div>
  );
}

export default Counter;
`
	if err := os.WriteFile(counterComponentPath, []byte(counterComponentContent), 0644); err != nil {
		return fmt.Errorf("failed to create Counter.jsx: %w", err)
	}

	// Create CSS for the counter component
	counterCSSPath := filepath.Join(projectPath, "src/components", "Counter.css")
	counterCSSContent := `.counter {
  text-align: center;
  margin: 2rem auto;
  padding: 1rem;
  max-width: 300px;
  border: 1px solid #ccc;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.counter-value {
  font-size: 3rem;
  font-weight: bold;
  margin: 1rem 0;
}

.counter-buttons {
  display: flex;
  justify-content: center;
  gap: 0.5rem;
}

.counter-buttons button {
  padding: 0.5rem 1rem;
  font-size: 1.25rem;
  border: none;
  border-radius: 4px;
  background-color: #0070f3;
  color: white;
  cursor: pointer;
}

.counter-buttons button:hover {
  background-color: #0060df;
}
`
	if err := os.WriteFile(counterCSSPath, []byte(counterCSSContent), 0644); err != nil {
		return fmt.Errorf("failed to create Counter.css: %w", err)
	}

	return nil
}

// setupReactMobXArchitecture sets up a React project with MobX
func setupReactMobXArchitecture(projectPath string) error {
	// Install MobX dependencies
	cmd := exec.Command("npm", "install", "mobx", "mobx-react-lite")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install MobX dependencies: %w", err)
	}

	// Create directories for MobX architecture
	dirs := []string{
		"src/components",
		"src/hooks",
		"src/utils",
		"src/assets",
		"src/stores",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	// Create a MobX store
	storePath := filepath.Join(projectPath, "src/stores", "counterStore.js")
	storeContent := `import { makeAutoObservable } from 'mobx';

class CounterStore {
  count = 0;

  constructor() {
    makeAutoObservable(this);
  }

  increment() {
    this.count += 1;
  }

  decrement() {
    this.count -= 1;
  }

  incrementByAmount(amount) {
    this.count += amount;
  }

  reset() {
    this.count = 0;
  }
}

// Create a singleton instance
const counterStore = new CounterStore();

export default counterStore;
`
	if err := os.WriteFile(storePath, []byte(storeContent), 0644); err != nil {
		return fmt.Errorf("failed to create counterStore.js: %w", err)
	}

	// Create a sample counter component using MobX
	counterComponentPath := filepath.Join(projectPath, "src/components", "Counter.jsx")
	counterComponentContent := `import React from 'react';
import { observer } from 'mobx-react-lite';
import counterStore from '../stores/counterStore';
import './Counter.css';

const Counter = observer(() => {
  return (
    <div className="counter">
      <h2>MobX Counter</h2>
      <div className="counter-value">{counterStore.count}</div>
      <div className="counter-buttons">
        <button onClick={() => counterStore.decrement()}>-</button>
        <button onClick={() => counterStore.increment()}>+</button>
        <button onClick={() => counterStore.incrementByAmount(5)}>+5</button>
      </div>
    </div>
  );
});

export default Counter;
`
	if err := os.WriteFile(counterComponentPath, []byte(counterComponentContent), 0644); err != nil {
		return fmt.Errorf("failed to create Counter.jsx: %w", err)
	}

	// Create CSS for the counter component
	counterCSSPath := filepath.Join(projectPath, "src/components", "Counter.css")
	counterCSSContent := `.counter {
  text-align: center;
  margin: 2rem auto;
  padding: 1rem;
  max-width: 300px;
  border: 1px solid #ccc;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.counter-value {
  font-size: 3rem;
  font-weight: bold;
  margin: 1rem 0;
}

.counter-buttons {
  display: flex;
  justify-content: center;
  gap: 0.5rem;
}

.counter-buttons button {
  padding: 0.5rem 1rem;
  font-size: 1.25rem;
  border: none;
  border-radius: 4px;
  background-color: #0070f3;
  color: white;
  cursor: pointer;
}

.counter-buttons button:hover {
  background-color: #0060df;
}
`
	if err := os.WriteFile(counterCSSPath, []byte(counterCSSContent), 0644); err != nil {
		return fmt.Errorf("failed to create Counter.css: %w", err)
	}

	return nil
}

// setupReactRecoilArchitecture sets up a React project with Recoil
func setupReactRecoilArchitecture(projectPath string) error {
	// Install Recoil
	cmd := exec.Command("npm", "install", "recoil")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Recoil: %w", err)
	}

	// Create directories for Recoil architecture
	dirs := []string{
		"src/components",
		"src/hooks",
		"src/utils",
		"src/assets",
		"src/atoms",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	// Create a Recoil atom
	atomPath := filepath.Join(projectPath, "src/atoms", "counterAtom.js")
	atomContent := `import { atom } from 'recoil';

export const counterState = atom({
  key: 'counterState', // unique ID
  default: 0, // default value
});
`
	if err := os.WriteFile(atomPath, []byte(atomContent), 0644); err != nil {
		return fmt.Errorf("failed to create counterAtom.js: %w", err)
	}

	// Update index.js to include Recoil provider
	indexPath := filepath.Join(projectPath, "src", "index.js")

	// Check if index.js exists
	_, err := os.Stat(indexPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to check index.js: %w", err)
	}

	// Create new content with Recoil provider
	newIndexContent := `import React from 'react';
import ReactDOM from 'react-dom/client';
import { RecoilRoot } from 'recoil';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <RecoilRoot>
      <App />
    </RecoilRoot>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
`

	// Write the updated content
	if err := os.WriteFile(indexPath, []byte(newIndexContent), 0644); err != nil {
		return fmt.Errorf("failed to update index.js: %w", err)
	}

	// Create a sample counter component using Recoil
	counterComponentPath := filepath.Join(projectPath, "src/components", "Counter.jsx")
	counterComponentContent := `import React from 'react';
import { useRecoilState } from 'recoil';
import { counterState } from '../atoms/counterAtom';
import './Counter.css';

function Counter() {
  const [count, setCount] = useRecoilState(counterState);

  const increment = () => setCount(count + 1);
  const decrement = () => setCount(count - 1);
  const incrementByAmount = (amount) => setCount(count + amount);

  return (
    <div className="counter">
      <h2>Recoil Counter</h2>
      <div className="counter-value">{count}</div>
      <div className="counter-buttons">
        <button onClick={decrement}>-</button>
        <button onClick={increment}>+</button>
        <button onClick={() => incrementByAmount(5)}>+5</button>
      </div>
    </div>
  );
}

export default Counter;
`
	if err := os.WriteFile(counterComponentPath, []byte(counterComponentContent), 0644); err != nil {
		return fmt.Errorf("failed to create Counter.jsx: %w", err)
	}

	// Create CSS for the counter component
	counterCSSPath := filepath.Join(projectPath, "src/components", "Counter.css")
	counterCSSContent := `.counter {
  text-align: center;
  margin: 2rem auto;
  padding: 1rem;
  max-width: 300px;
  border: 1px solid #ccc;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.counter-value {
  font-size: 3rem;
  font-weight: bold;
  margin: 1rem 0;
}

.counter-buttons {
  display: flex;
  justify-content: center;
  gap: 0.5rem;
}

.counter-buttons button {
  padding: 0.5rem 1rem;
  font-size: 1.25rem;
  border: none;
  border-radius: 4px;
  background-color: #0070f3;
  color: white;
  cursor: pointer;
}

.counter-buttons button:hover {
  background-color: #0060df;
}
`
	if err := os.WriteFile(counterCSSPath, []byte(counterCSSContent), 0644); err != nil {
		return fmt.Errorf("failed to create Counter.css: %w", err)
	}

	return nil
}

// getReactArchitectureName returns a human-readable name for the architecture
func getReactArchitectureName(stateManagement string) string {
	switch strings.ToLower(stateManagement) {
	case "redux":
		return "Redux"
	case "context":
		return "Context API"
	case "mobx":
		return "MobX"
	case "recoil":
		return "Recoil"
	default:
		return "basic"
	}
}
