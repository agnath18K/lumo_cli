package create

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// addDependenciesToPubspec adds dependencies to a Flutter project's pubspec.yaml
func addDependenciesToPubspec(pubspecPath string, dependencies []string) error {
	// Read the pubspec.yaml file
	content, err := os.ReadFile(pubspecPath)
	if err != nil {
		return fmt.Errorf("failed to read pubspec.yaml: %w", err)
	}

	// Convert to string for easier manipulation
	pubspecContent := string(content)

	// Find the dependencies section
	dependenciesIndex := strings.Index(pubspecContent, "dependencies:")
	if dependenciesIndex == -1 {
		return fmt.Errorf("could not find dependencies section in pubspec.yaml")
	}

	// Find the next section after dependencies
	nextSectionIndex := strings.Index(pubspecContent[dependenciesIndex:], "\n\n")
	if nextSectionIndex == -1 {
		// If there's no next section, use the end of the file
		nextSectionIndex = len(pubspecContent) - dependenciesIndex
	} else {
		nextSectionIndex += dependenciesIndex
	}

	// Extract the dependencies section
	dependenciesSection := pubspecContent[dependenciesIndex:nextSectionIndex]

	// Add new dependencies
	newDependenciesSection := dependenciesSection
	for _, dep := range dependencies {
		// Check if the dependency is already in the file
		if !strings.Contains(newDependenciesSection, strings.Split(dep, ":")[0]) {
			newDependenciesSection += fmt.Sprintf("\n  %s", dep)
		}
	}

	// Replace the old dependencies section with the new one
	newPubspecContent := pubspecContent[:dependenciesIndex] + newDependenciesSection + pubspecContent[nextSectionIndex:]

	// Write the updated content back to the file
	if err := os.WriteFile(pubspecPath, []byte(newPubspecContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated pubspec.yaml: %w", err)
	}

	// Run flutter pub get to install the dependencies
	cmd := exec.Command("flutter", "pub", "get")
	cmd.Dir = filepath.Dir(pubspecPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run 'flutter pub get': %w", err)
	}

	return nil
}

// createSampleBlocFiles creates sample BLoC files for a Flutter project
func createSampleBlocFiles(projectPath string) error {
	// Create a sample counter bloc
	counterBlocPath := filepath.Join(projectPath, "lib/blocs/counter_bloc.dart")
	counterBlocContent := `import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';

// Events
abstract class CounterEvent extends Equatable {
  const CounterEvent();

  @override
  List<Object> get props => [];
}

class IncrementEvent extends CounterEvent {}

class DecrementEvent extends CounterEvent {}

// States
class CounterState extends Equatable {
  final int count;

  const CounterState(this.count);

  @override
  List<Object> get props => [count];
}

// BLoC
class CounterBloc extends Bloc<CounterEvent, CounterState> {
  CounterBloc() : super(const CounterState(0)) {
    on<IncrementEvent>((event, emit) {
      emit(CounterState(state.count + 1));
    });

    on<DecrementEvent>((event, emit) {
      emit(CounterState(state.count - 1));
    });
  }
}
`
	if err := os.WriteFile(counterBlocPath, []byte(counterBlocContent), 0644); err != nil {
		return fmt.Errorf("failed to create counter_bloc.dart: %w", err)
	}

	// Create a sample counter screen
	counterScreenPath := filepath.Join(projectPath, "lib/screens/counter_screen.dart")
	counterScreenContent := `import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../blocs/counter_bloc.dart';

class CounterScreen extends StatelessWidget {
  const CounterScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Counter Example'),
      ),
      body: Center(
        child: BlocBuilder<CounterBloc, CounterState>(
          builder: (context, state) {
            return Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: <Widget>[
                const Text(
                  'You have pushed the button this many times:',
                ),
                Text(
                  '${state.count}',
                  style: Theme.of(context).textTheme.headlineMedium,
                ),
              ],
            );
          },
        ),
      ),
      floatingActionButton: Column(
        mainAxisAlignment: MainAxisAlignment.end,
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          FloatingActionButton(
            onPressed: () => context.read<CounterBloc>().add(IncrementEvent()),
            tooltip: 'Increment',
            child: const Icon(Icons.add),
          ),
          const SizedBox(height: 8),
          FloatingActionButton(
            onPressed: () => context.read<CounterBloc>().add(DecrementEvent()),
            tooltip: 'Decrement',
            child: const Icon(Icons.remove),
          ),
        ],
      ),
    );
  }
}
`
	if err := os.WriteFile(counterScreenPath, []byte(counterScreenContent), 0644); err != nil {
		return fmt.Errorf("failed to create counter_screen.dart: %w", err)
	}

	return nil
}

// updateMainDartForBloc updates the main.dart file to use BLoC
func updateMainDartForBloc(projectPath string) error {
	mainDartPath := filepath.Join(projectPath, "lib/main.dart")
	mainDartContent := `import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'blocs/counter_bloc.dart';
import 'screens/counter_screen.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter BLoC Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: BlocProvider(
        create: (context) => CounterBloc(),
        child: const CounterScreen(),
      ),
    );
  }
}
`
	if err := os.WriteFile(mainDartPath, []byte(mainDartContent), 0644); err != nil {
		return fmt.Errorf("failed to update main.dart: %w", err)
	}

	return nil
}

// createSampleProviderFiles creates sample Provider files for a Flutter project
func createSampleProviderFiles(projectPath string) error {
	// Create a sample counter provider
	counterProviderPath := filepath.Join(projectPath, "lib/providers/counter_provider.dart")
	counterProviderContent := `import 'package:flutter/foundation.dart';

class CounterProvider with ChangeNotifier {
  int _count = 0;

  int get count => _count;

  void increment() {
    _count++;
    notifyListeners();
  }

  void decrement() {
    _count--;
    notifyListeners();
  }
}
`
	if err := os.WriteFile(counterProviderPath, []byte(counterProviderContent), 0644); err != nil {
		return fmt.Errorf("failed to create counter_provider.dart: %w", err)
	}

	// Create a sample counter screen
	counterScreenPath := filepath.Join(projectPath, "lib/screens/counter_screen.dart")
	counterScreenContent := `import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/counter_provider.dart';

class CounterScreen extends StatelessWidget {
  const CounterScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Counter Example'),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            const Text(
              'You have pushed the button this many times:',
            ),
            Consumer<CounterProvider>(
              builder: (context, counter, child) {
                return Text(
                  '${counter.count}',
                  style: Theme.of(context).textTheme.headlineMedium,
                );
              },
            ),
          ],
        ),
      ),
      floatingActionButton: Column(
        mainAxisAlignment: MainAxisAlignment.end,
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          FloatingActionButton(
            onPressed: () => Provider.of<CounterProvider>(context, listen: false).increment(),
            tooltip: 'Increment',
            child: const Icon(Icons.add),
          ),
          const SizedBox(height: 8),
          FloatingActionButton(
            onPressed: () => Provider.of<CounterProvider>(context, listen: false).decrement(),
            tooltip: 'Decrement',
            child: const Icon(Icons.remove),
          ),
        ],
      ),
    );
  }
}
`
	if err := os.WriteFile(counterScreenPath, []byte(counterScreenContent), 0644); err != nil {
		return fmt.Errorf("failed to create counter_screen.dart: %w", err)
	}

	return nil
}

// updateMainDartForProvider updates the main.dart file to use Provider
func updateMainDartForProvider(projectPath string) error {
	mainDartPath := filepath.Join(projectPath, "lib/main.dart")
	mainDartContent := `import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'providers/counter_provider.dart';
import 'screens/counter_screen.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (context) => CounterProvider(),
      child: MaterialApp(
        title: 'Flutter Provider Demo',
        theme: ThemeData(
          primarySwatch: Colors.blue,
        ),
        home: const CounterScreen(),
      ),
    );
  }
}
`
	if err := os.WriteFile(mainDartPath, []byte(mainDartContent), 0644); err != nil {
		return fmt.Errorf("failed to update main.dart: %w", err)
	}

	return nil
}
