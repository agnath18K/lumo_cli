package create

import (
	"fmt"
	"os"
	"path/filepath"
)

// createSampleRiverpodFiles creates sample Riverpod files for a Flutter project
func createSampleRiverpodFiles(projectPath string) error {
	// Create a sample counter provider
	counterProviderPath := filepath.Join(projectPath, "lib/providers/counter_provider.dart")
	counterProviderContent := `import 'package:flutter_riverpod/flutter_riverpod.dart';

final counterProvider = StateNotifierProvider<CounterNotifier, int>((ref) {
  return CounterNotifier();
});

class CounterNotifier extends StateNotifier<int> {
  CounterNotifier() : super(0);
  
  void increment() => state = state + 1;
  void decrement() => state = state - 1;
}
`
	if err := os.WriteFile(counterProviderPath, []byte(counterProviderContent), 0644); err != nil {
		return fmt.Errorf("failed to create counter_provider.dart: %w", err)
	}
	
	// Create a sample counter screen
	counterScreenPath := filepath.Join(projectPath, "lib/screens/counter_screen.dart")
	counterScreenContent := `import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/counter_provider.dart';

class CounterScreen extends ConsumerWidget {
  const CounterScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final count = ref.watch(counterProvider);
    
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
            Text(
              '$count',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
          ],
        ),
      ),
      floatingActionButton: Column(
        mainAxisAlignment: MainAxisAlignment.end,
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          FloatingActionButton(
            onPressed: () => ref.read(counterProvider.notifier).increment(),
            tooltip: 'Increment',
            child: const Icon(Icons.add),
          ),
          const SizedBox(height: 8),
          FloatingActionButton(
            onPressed: () => ref.read(counterProvider.notifier).decrement(),
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

// updateMainDartForRiverpod updates the main.dart file to use Riverpod
func updateMainDartForRiverpod(projectPath string) error {
	mainDartPath := filepath.Join(projectPath, "lib/main.dart")
	mainDartContent := `import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'screens/counter_screen.dart';

void main() {
  runApp(const ProviderScope(child: MyApp()));
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Riverpod Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: const CounterScreen(),
    );
  }
}
`
	if err := os.WriteFile(mainDartPath, []byte(mainDartContent), 0644); err != nil {
		return fmt.Errorf("failed to update main.dart: %w", err)
	}
	
	return nil
}

// createSampleMVVMFiles creates sample MVVM files for a Flutter project
func createSampleMVVMFiles(projectPath string) error {
	// Create a sample counter model
	counterModelPath := filepath.Join(projectPath, "lib/models/counter_model.dart")
	counterModelContent := `class CounterModel {
  int count;
  
  CounterModel({this.count = 0});
}
`
	if err := os.WriteFile(counterModelPath, []byte(counterModelContent), 0644); err != nil {
		return fmt.Errorf("failed to create counter_model.dart: %w", err)
	}
	
	// Create a sample counter view model
	counterViewModelPath := filepath.Join(projectPath, "lib/viewmodels/counter_viewmodel.dart")
	counterViewModelContent := `import 'package:flutter/foundation.dart';
import '../models/counter_model.dart';

class CounterViewModel with ChangeNotifier {
  final CounterModel _model = CounterModel();
  
  int get count => _model.count;
  
  void increment() {
    _model.count++;
    notifyListeners();
  }
  
  void decrement() {
    _model.count--;
    notifyListeners();
  }
}
`
	if err := os.WriteFile(counterViewModelPath, []byte(counterViewModelContent), 0644); err != nil {
		return fmt.Errorf("failed to create counter_viewmodel.dart: %w", err)
	}
	
	// Create a sample counter view
	counterViewPath := filepath.Join(projectPath, "lib/views/counter_view.dart")
	counterViewContent := `import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../viewmodels/counter_viewmodel.dart';

class CounterView extends StatelessWidget {
  const CounterView({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final viewModel = Provider.of<CounterViewModel>(context);
    
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
            Text(
              '${viewModel.count}',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
          ],
        ),
      ),
      floatingActionButton: Column(
        mainAxisAlignment: MainAxisAlignment.end,
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          FloatingActionButton(
            onPressed: viewModel.increment,
            tooltip: 'Increment',
            child: const Icon(Icons.add),
          ),
          const SizedBox(height: 8),
          FloatingActionButton(
            onPressed: viewModel.decrement,
            tooltip: 'Decrement',
            child: const Icon(Icons.remove),
          ),
        ],
      ),
    );
  }
}
`
	if err := os.WriteFile(counterViewPath, []byte(counterViewContent), 0644); err != nil {
		return fmt.Errorf("failed to create counter_view.dart: %w", err)
	}
	
	return nil
}

// updateMainDartForMVVM updates the main.dart file to use MVVM
func updateMainDartForMVVM(projectPath string) error {
	mainDartPath := filepath.Join(projectPath, "lib/main.dart")
	mainDartContent := `import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'viewmodels/counter_viewmodel.dart';
import 'views/counter_view.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (context) => CounterViewModel(),
      child: MaterialApp(
        title: 'Flutter MVVM Demo',
        theme: ThemeData(
          primarySwatch: Colors.blue,
        ),
        home: const CounterView(),
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
