# Go-DecisionTree-ML

A from-scratch implementation of a **Decision Tree** machine learning algorithm using pure Go. 

Instead of relying on external ML frameworks, this project builds the foundational mathematics of decision trees directly into the code. It dynamically reads categorical data, calculates **Entropy** and **Information Gain** to find the optimal splits, and recursively constructs the tree to make predictions.

## Features
* **Zero External ML Dependencies:** Built using only Go standard libraries (`math`, `encoding/csv`, `fmt`, etc.).
* **Dynamic Tree Building:** Automatically calculates Information Gain to select the best root and child nodes.
* **Interactive Prediction CLI:** Includes a command-line interface that traverses the generated tree based on user input to return a classification.

## How It Works
The algorithm follows the ID3 (Iterative Dichotomiser 3) logic:
1. Calculates the Entropy of the dataset.
2. Iterates through attributes to calculate the Information Gain for each.
3. Splits the dataset based on the attribute with the highest Information Gain.
4. Recursively builds a tree structure (using Go `map` interfaces) until all subsets are pure or no attributes remain.

## Usage

### 1. Data Preparation
The program expects a categorical CSV file named `dataset.csv` in the root directory. The last column **must** be the target class label.

**Example `dataset.csv` format:**
```csv
Outlook,Temperature,Humidity,Wind,PlayTennis
Sunny,Hot,High,Weak,No
Sunny,Hot,High,Strong,No
Overcast,Hot,High,Weak,Yes
Rain,Mild,High,Weak,Yes
Rain,Cool,Normal,Weak,Yes
```

### 2. Running the Model
Ensure you have Go installed on your machine, then run:

```bash
go run decision_tree.go
```

### 3. Example Output
The program will output the Information Gain calculations, print the tree structure, and prompt you for input:

```plaintext
Information gain for Outlook: 0.2467
Information gain for Temperature: 0.0292
Information gain for Humidity: 0.1518

Prediction Tree:
Outlook: Sunny ->
  Humidity: High ->
    Prediction: No
  Humidity: Normal ->
    Prediction: Yes

Enter attribute values for prediction (press Enter to quit):
Outlook: Sunny
Temperature: Mild
Humidity: Normal
Wind: Weak
Prediction: Yes
```
