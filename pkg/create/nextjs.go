package create

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// generateNextJSProject creates a new Next.js project
func generateNextJSProject(stateManagement string, options map[string]string) (string, error) {
	// Get project name from options or use a default
	projectName := options["name"]
	if projectName == "" {
		projectName = "my-nextjs-app"
	}

	// Check if Node.js is installed
	if err := checkNodeInstalled(); err != nil {
		return "", err
	}

	// Create the project using create-next-app
	if err := createBaseNextJSProject(projectName); err != nil {
		return "", err
	}

	// Set up the project structure based on state management
	switch strings.ToLower(stateManagement) {
	case "redux":
		if err := setupReduxArchitecture(projectName); err != nil {
			return "", err
		}
	case "context":
		if err := setupContextAPIArchitecture(projectName); err != nil {
			return "", err
		}
	case "zustand":
		if err := setupZustandArchitecture(projectName); err != nil {
			return "", err
		}
	default:
		// Default to a basic structure without specific state management
		if err := setupBasicNextJSArchitecture(projectName); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("✅ Next.js project '%s' created successfully with %s architecture!",
		projectName,
		getNextJSArchitectureName(stateManagement)), nil
}

// checkNodeInstalled verifies that Node.js is installed
func checkNodeInstalled() error {
	cmd := exec.Command("node", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Node.js is not installed or not in PATH. Please install Node.js first: https://nodejs.org/")
	}
	return nil
}

// createBaseNextJSProject creates a new Next.js project using create-next-app
func createBaseNextJSProject(name string) error {
	// Use npx to run create-next-app without installing it globally
	cmd := exec.Command("npx", "create-next-app@latest", name, "--use-npm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// setupBasicNextJSArchitecture sets up a basic Next.js project structure
func setupBasicNextJSArchitecture(projectPath string) error {
	// Create additional directories for a clean architecture
	dirs := []string{
		"components",
		"lib",
		"utils",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	// Create a sample utility function
	utilsPath := filepath.Join(projectPath, "utils", "helpers.js")
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
	buttonComponentPath := filepath.Join(projectPath, "components", "Button.jsx")
	buttonComponentContent := `import React from 'react';

/**
 * Button component with variants
 */
export default function Button({ children, variant = 'primary', onClick }) {
  const baseStyles = 'px-4 py-2 rounded font-medium focus:outline-none focus:ring-2 focus:ring-offset-2';

  const variantStyles = {
    primary: 'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500',
    secondary: 'bg-gray-200 text-gray-800 hover:bg-gray-300 focus:ring-gray-500',
    danger: 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500',
  };

  const styles = baseStyles + ' ' + (variantStyles[variant] || variantStyles.primary);

  return (
    <button className={styles} onClick={onClick}>
      {children}
    </button>
  );
}
`
	if err := os.WriteFile(buttonComponentPath, []byte(buttonComponentContent), 0644); err != nil {
		return fmt.Errorf("failed to create Button.jsx: %w", err)
	}

	return nil
}

// setupReduxArchitecture sets up a Next.js project with Redux
func setupReduxArchitecture(projectPath string) error {
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
		"components",
		"lib",
		"utils",
		"store",
		"store/slices",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	// Create Redux store
	storePath := filepath.Join(projectPath, "store", "index.js")
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
	slicePath := filepath.Join(projectPath, "store", "slices", "counterSlice.js")
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

	// Update _app.js to include Redux provider
	appPath := filepath.Join(projectPath, "pages", "_app.js")

	// Check if _app.js exists
	_, err := os.Stat(appPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to check _app.js: %w", err)
	}

	// Create new content with Redux provider
	newAppContent := `import { Provider } from 'react-redux';
import { store } from '../store';
import '../styles/globals.css';

function MyApp({ Component, pageProps }) {
  return (
    <Provider store={store}>
      <Component {...pageProps} />
    </Provider>
  );
}

export default MyApp;
`

	// Write the updated content
	if err := os.WriteFile(appPath, []byte(newAppContent), 0644); err != nil {
		return fmt.Errorf("failed to update _app.js: %w", err)
	}

	return nil
}

// setupContextAPIArchitecture sets up a Next.js project with Context API
func setupContextAPIArchitecture(projectPath string) error {
	// Create directories for Context API architecture
	dirs := []string{
		"components",
		"lib",
		"utils",
		"contexts",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	// Create a sample context
	contextPath := filepath.Join(projectPath, "contexts", "CounterContext.jsx")
	contextContent := `import React, { createContext, useContext, useState } from 'react';

// Create the context
const CounterContext = createContext();

// Create a provider component
export function CounterProvider({ children }) {
  const [count, setCount] = useState(0);

  const increment = () => setCount(count + 1);
  const decrement = () => setCount(count - 1);
  const reset = () => setCount(0);

  const value = {
    count,
    increment,
    decrement,
    reset,
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

	// Update _app.js to include Context provider
	appPath := filepath.Join(projectPath, "pages", "_app.js")

	// Check if _app.js exists
	_, err := os.Stat(appPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to check _app.js: %w", err)
	}

	// Create new content with Context provider
	newAppContent := `import { CounterProvider } from '../contexts/CounterContext';
import '../styles/globals.css';

function MyApp({ Component, pageProps }) {
  return (
    <CounterProvider>
      <Component {...pageProps} />
    </CounterProvider>
  );
}

export default MyApp;
`

	// Write the updated content
	if err := os.WriteFile(appPath, []byte(newAppContent), 0644); err != nil {
		return fmt.Errorf("failed to update _app.js: %w", err)
	}

	return nil
}

// setupZustandArchitecture sets up a Next.js project with Zustand
func setupZustandArchitecture(projectPath string) error {
	// Install Zustand
	cmd := exec.Command("npm", "install", "zustand")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Zustand: %w", err)
	}

	// Create directories for Zustand architecture
	dirs := []string{
		"components",
		"lib",
		"utils",
		"store",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	// Create a sample Zustand store
	storePath := filepath.Join(projectPath, "store", "useCounterStore.js")
	storeContent := `import create from 'zustand';

// Create a store with Zustand
const useCounterStore = create((set) => ({
  count: 0,
  increment: () => set((state) => ({ count: state.count + 1 })),
  decrement: () => set((state) => ({ count: state.count - 1 })),
  reset: () => set({ count: 0 }),
  incrementByAmount: (amount) => set((state) => ({ count: state.count + amount })),
}));

export default useCounterStore;
`
	if err := os.WriteFile(storePath, []byte(storeContent), 0644); err != nil {
		return fmt.Errorf("failed to create useCounterStore.js: %w", err)
	}

	return nil
}

// getNextJSArchitectureName returns a human-readable name for the architecture
func getNextJSArchitectureName(stateManagement string) string {
	switch strings.ToLower(stateManagement) {
	case "redux":
		return "Redux"
	case "context":
		return "Context API"
	case "zustand":
		return "Zustand"
	default:
		return "basic"
	}
}
