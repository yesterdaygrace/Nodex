# PRD: Sentiment Analysis Dicoding Submission

## 1. Project Overview

- **Project Name:** Indonesian Sentiment Analysis
- **Project Type:** Machine Learning / NLP
- **Primary Language:** Python
- **Target:** Dicoding Sentiment Analysis Submission
- **Dataset Target:** ≥10,000 independently scraped samples
- **Classes:** Negative, Neutral, Positive
- **Target Model:** Deep Learning
- **Target Accuracy:** >92% training and testing accuracy
- **Minimum Acceptance:** ≥85% testing accuracy

### Objective

Build a reproducible Indonesian sentiment-analysis pipeline that:

1. Independently scrapes raw opinion/review data.
2. Cleans and labels the data.
3. Performs feature extraction/tokenization.
4. Trains and compares at least 3 different training schemes.
5. Achieves the strongest possible testing performance, targeting >92%.
6. Provides notebook-based inference producing categorical sentiment.
7. Packages everything into one `.zip` submission.

---

## 2. Submission Requirements

| Requirement | Target |
|---|---|
| Programming language | Python |
| Independent scraping | Required |
| Minimum dataset | 3,000 samples |
| Recommended dataset | 10,000+ samples |
| Feature extraction | Required |
| Data labeling | Required |
| ML training | Required |
| Minimum testing accuracy | ≥85% |
| Recommended testing accuracy | >92% |
| Dataset classes | Recommended ≥3 |
| Training experiments | Recommended 3 |
| Inference | Recommended |
| Submission format | `.zip` |
| Executed notebook | Required |

---

## 3. Recommended Final Target

The implementation should aim for the Rating 5 configuration.

- 10,000–20,000+ scraped samples
- 3 sentiment classes:
  - Negative
  - Neutral
  - Positive
- Independent scraping using Python
- Automatic/semiautomatic labeling from rating signals
- Deep-learning model
- 3 experimental training schemes
- At least 2 dimensions changed between experiments
- Best model:
  - Training accuracy >92%
  - Testing accuracy >92%
- Inference cell inside `.ipynb`
- Visible inference output
- Clean `requirements.txt`
- Complete scraping source code
- No pre-existing open-source sentiment dataset

---

## 4. System Architecture

```text
                    ┌──────────────────────┐
                    │     Data Source      │
                    │ Google Play / Other  │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │   Python Scraper     │
                    │ requests / API /     │
                    │ permitted scraping   │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │    Raw Dataset       │
                    │ raw_reviews.csv      │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │ Data Cleaning        │
                    │ - duplicate removal │
                    │ - normalization     │
                    │ - empty removal     │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │     Labeling         │
                    │ Negative / Neutral   │
                    │ Positive             │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │ Train / Validation   │
                    │ / Test Split         │
                    └──────────┬───────────┘
                               │
                ┌──────────────┼──────────────┐
                ▼              ▼              ▼
          Experiment 1   Experiment 2   Experiment 3
                │              │              │
                └──────────────┼──────────────┘
                               ▼
                    ┌──────────────────────┐
                    │ Model Evaluation     │
                    │ Accuracy / Report    │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │   Best Model         │
                    └──────────┬───────────┘
                               │
                               ▼
                    ┌──────────────────────┐
                    │     Inference        │
                    │ Text → Sentiment     │
                    └──────────────────────┘
```

---

## 5. Data Source Strategy

### Recommended Source: Google Play Store Reviews

This is the strongest practical option because it provides:

- Large quantities of Indonesian-language reviews.
- Natural user-generated text.
- Star ratings that can be converted into sentiment labels.
- Clear sentiment signals.
- A practical Python collection workflow.
- A strong explanation for why the dataset was collected.

### Recommended Raw Schema

```text
review_id
review_text
rating
date
app_id
app_name
```

---

## 6. Dataset Requirements

### Minimum

```text
3,000 samples
```

### Recommended

```text
10,000+ samples
```

### Ideal

```text
15,000–30,000 raw samples
```

Do not assume that 10,000 scraped rows equals 10,000 usable rows.

Example:

```text
Raw scraping
    ↓
15,000
    ↓
duplicate removal
    ↓
13,500
    ↓
empty/invalid removal
    ↓
12,800
    ↓
usable dataset
```

Target at least **10,000 clean labeled samples**.

---

## 7. Sentiment Labeling

Use rating-based weak supervision.

| Rating | Sentiment |
|---:|---|
| 1 | Negative |
| 2 | Negative |
| 3 | Neutral |
| 4 | Positive |
| 5 | Positive |

```python
if rating <= 2:
    label = "negative"
elif rating == 3:
    label = "neutral"
else:
    label = "positive"
```

Document this explicitly as **rating-derived labeling**.

---

## 8. Dataset Quality Pipeline

```text
Raw reviews
    ↓
Remove empty reviews
    ↓
Remove duplicates
    ↓
Normalize whitespace
    ↓
Normalize URLs
    ↓
Normalize mentions
    ↓
Normalize repeated characters
    ↓
Handle emojis
    ↓
Remove invalid rows
    ↓
Label sentiment
    ↓
Check class distribution
    ↓
Final dataset
```

Do not aggressively remove punctuation or emojis without testing because they may contain sentiment information.

---

## 9. Dataset Structure

Recommended final CSV:

```text
dataset/
└── sentiment_dataset.csv
```

Columns:

```text
id
text
rating
label
source
created_at
```

Example:

| id | text | rating | label | source |
|---|---|---:|---|---|
| 1 | aplikasinya sangat bagus | 5 | positive | google_play |
| 2 | sering error dan lambat | 1 | negative | google_play |
| 3 | cukup baik tapi masih ada bug | 3 | neutral | google_play |

---

## 10. Feature Extraction / Text Representation

### Experiment A: TF-IDF + Classical ML

```text
Text
 ↓
TF-IDF
 ↓
Linear SVM
 ↓
Sentiment
```

Purpose: establish a strong classical baseline.

### Experiment B: Embedding + Deep Learning

```text
Text
 ↓
Tokenizer
 ↓
Padding
 ↓
Embedding
 ↓
BiLSTM
 ↓
Dense
 ↓
Softmax
 ↓
Sentiment
```

Purpose: introduce deep learning and sequence modeling.

### Experiment C: CNN + BiLSTM

```text
Text
 ↓
Tokenizer
 ↓
Padding
 ↓
Embedding
 ↓
Conv1D
 ↓
BiLSTM
 ↓
Dense
 ↓
Softmax
 ↓
Sentiment
```

Purpose: compare a second deep-learning architecture.

---

## 11. Training Experiments

At least 3 experiments should be documented.

| Experiment | Feature | Model | Split |
|---|---|---|---|
| E1 | TF-IDF | Linear SVM | 80/20 |
| E2 | Embedding | BiLSTM | 80/10/10 |
| E3 | Embedding | CNN + BiLSTM | 80/10/10 |

This provides multiple combinations through changes in:

- Feature extraction
- Training algorithm
- Model architecture
- Data split

---

## 12. Deep Learning Architecture

### Recommended Model

Embedding → BiLSTM → Dropout → Dense

```text
Input Text
     │
     ▼
Tokenizer
     │
     ▼
Padding
     │
     ▼
Embedding
     │
     ▼
Bidirectional LSTM
     │
     ▼
Dropout
     │
     ▼
Dense
     │
     ▼
Softmax
     │
     ▼
3 Classes
```

Output:

```text
[negative_probability,
 neutral_probability,
 positive_probability]
```

Final output:

```text
positive
```

---

## 13. Model Training

Recommended starting configuration:

```text
Optimizer: Adam
Loss: Sparse Categorical Crossentropy
Metric: Accuracy
Batch Size: 32 / 64
Epochs: 10–20
Early Stopping: Enabled
```

Use validation data to detect overfitting.

Do not optimize only for training accuracy.

A model with:

```text
Training = 99%
Testing  = 71%
```

is overfitting and does not satisfy the project objective.

---

## 14. Accuracy Targets

### Absolute minimum

```text
Testing Accuracy >= 85%
```

### Rating 4

All main criteria plus at least 3 recommendations.

### Rating 5

All main criteria plus all recommendations.

### Ideal

```text
Training Accuracy > 92%
Testing Accuracy > 92%
```

Also calculate:

- Precision
- Recall
- F1-score
- Confusion matrix

---

## 15. Class Distribution

Always inspect class distribution before training.

Example:

```text
Positive    5,200
Negative    4,100
Neutral     1,200
```

Severe imbalance can make accuracy misleading.

Consider balancing strategies if necessary.

---

## 16. Train / Validation / Test Split

Recommended:

```text
80% Training
10% Validation
10% Testing
```

For 10,000 samples:

```text
Training      8,000
Validation    1,000
Testing       1,000
```

Use stratification where applicable.

---

## 17. Data Leakage Prevention

The pipeline must follow:

```text
RAW DATA
   ↓
CLEAN
   ↓
SPLIT
   ↓
FIT TOKENIZER/VECTORIZER ON TRAINING ONLY
   ↓
TRANSFORM VALIDATION
   ↓
TRANSFORM TEST
   ↓
TRAIN MODEL
```

Never fit TF-IDF or tokenization vocabulary on the complete dataset before splitting.

---

## 18. Inference Requirement

The notebook must contain an inference section if recommendation #6 is implemented.

Example:

```text
Input:
"Aplikasinya sekarang jauh lebih cepat dan mudah digunakan."

Output:
Positive
```

```text
Input:
"Setiap dibuka selalu crash, sangat mengecewakan."

Output:
Negative
```

```text
Input:
"Fiturnya lumayan, tapi masih ada beberapa kekurangan."

Output:
Neutral
```

The categorical output must be visible in the executed notebook.

---

## 19. Notebook Structure

Recommended `training.ipynb`:

```text
01. Project Introduction
02. Import Libraries
03. Configuration
04. Load Dataset
05. Dataset Exploration
06. Data Cleaning
07. Label Distribution
08. Train/Test Split

09. Experiment 1
    - TF-IDF
    - SVM
    - Evaluation

10. Experiment 2
    - Tokenization
    - Embedding
    - BiLSTM
    - Evaluation

11. Experiment 3
    - Tokenization
    - Embedding
    - CNN + BiLSTM
    - Evaluation

12. Experiment Comparison
13. Best Model Selection
14. Final Evaluation
15. Confusion Matrix
16. Classification Report
17. Inference
18. Conclusion
```

---

## 20. Scraping Code

Recommended:

```text
scrape.py
```

Responsibilities:

```text
Initialize scraper
    ↓
Collect reviews
    ↓
Handle pagination
    ↓
Normalize raw records
    ↓
Remove duplicates
    ↓
Save CSV
```

Do not put model training inside `scrape.py`.

---

## 21. Project Folder Structure

```text
sentiment-analysis/
│
├── README.md
├── requirements.txt
│
├── scraping/
│   └── scrape.py
│
├── notebooks/
│   └── training.ipynb
│
├── dataset/
│   └── sentiment_dataset.csv
│
├── src/
│   ├── preprocessing.py
│   ├── labeling.py
│   ├── evaluation.py
│   └── inference.py
│
└── results/
    ├── confusion_matrix.png
    ├── training_history.png
    └── experiment_results.csv
```

For the final Dicoding ZIP, simplify if extra files do not materially improve reproducibility.

---

## 22. requirements.txt

Two supported approaches:

### pip freeze

```bash
pip freeze > requirements.txt
```

### pipreqs

```bash
pipreqs .
```

Recommended: use `pipreqs`, then manually verify the resulting dependency list.

---

## 23. Recommended Libraries

Potential stack:

```text
Python
pandas
numpy
scikit-learn
tensorflow
matplotlib
seaborn
nltk
Sastrawi
requests
beautifulsoup4
```

Only include packages actually required by the project.

---

## 24. Evaluation

For every experiment record:

```text
Experiment
Model
Feature extraction
Train accuracy
Validation accuracy
Test accuracy
Precision
Recall
F1-score
Training time
```

Example:

| Model | Train | Test | F1 |
|---|---:|---:|---:|
| TF-IDF + SVM | 94.2% | 88.7% | 0.88 |
| BiLSTM | 95.4% | 91.2% | 0.91 |
| CNN + BiLSTM | 97.1% | 93.4% | 0.93 |

The third model becomes the candidate final model if the actual measured results support it.

---

## 25. Confusion Matrix

Generate a confusion matrix to identify class-level errors.

```text
                 Predicted
              Neg  Neu  Pos
Actual Neg     850  90   60
       Neu      80  820 100
       Pos      40  70  890
```

Accuracy should not be the only evidence of model quality.

---

## 26. Reproducibility

Set random seeds for:

- Python
- NumPy
- TensorFlow

Document:

```text
Python version
TensorFlow version
Dataset size
Random seed
Training configuration
```

---

## 27. Data Provenance

README must explicitly document:

```text
Dataset was collected independently using Python.

Source:
Google Play Store reviews

Collection method:
Python scraper

Collection date:
YYYY-MM-DD

Raw sample count:
XX,XXX

Final sample count:
XX,XXX

Labeling:
1–2 stars → Negative
3 stars → Neutral
4–5 stars → Positive
```

Do not claim independently scraped data if an existing Kaggle/open-source sentiment dataset was used.

---

## 28. README Requirements

README should contain:

```text
1. Project Overview
2. Objective
3. Dataset Source
4. Scraping Method
5. Dataset Statistics
6. Labeling Method
7. Preprocessing
8. Feature Extraction
9. Training Experiments
10. Model Architecture
11. Evaluation Results
12. Inference Example
13. Installation
14. Running Instructions
15. Project Structure
16. Reproducibility
```

---

## 29. Final ZIP Structure

```text
sentiment-analysis.zip
│
└── sentiment-analysis/
    │
    ├── README.md
    ├── requirements.txt
    │
    ├── scrape.py
    │
    ├── training.ipynb
    │
    └── sentiment_dataset.csv
```

Required artifacts:

- `training.ipynb`
- `scrape.py` or scraping notebook
- `requirements.txt`
- `sentiment_dataset.csv`

The notebook must already be executed so its outputs are visible.

---

## 30. Submission Checklist

### Data

- [ ] Data scraped independently
- [ ] Scraping code included
- [ ] At least 3,000 samples
- [ ] Target ≥10,000 samples
- [ ] Dataset saved as CSV/JSON
- [ ] Dataset is not an existing open-source dataset
- [ ] Duplicate data removed
- [ ] Empty/invalid rows removed
- [ ] Three sentiment classes available
- [ ] Class distribution documented

### Machine Learning

- [ ] Feature extraction implemented
- [ ] Labeling implemented
- [ ] Training algorithm implemented
- [ ] Three experiments performed
- [ ] At least two experimental dimensions differ
- [ ] Testing accuracy ≥85%
- [ ] Target >92% train/test
- [ ] Deep learning model included
- [ ] Precision calculated
- [ ] Recall calculated
- [ ] F1-score calculated
- [ ] Confusion matrix generated

### Inference

- [ ] Inference implemented
- [ ] Input text accepted
- [ ] Output is categorical
- [ ] Negative output demonstrated
- [ ] Neutral output demonstrated
- [ ] Positive output demonstrated
- [ ] Output visible in notebook

### Submission

- [ ] `training.ipynb`
- [ ] `scrape.py` or scraping notebook
- [ ] `requirements.txt`
- [ ] `sentiment_dataset.csv`
- [ ] Everything inside one ZIP
- [ ] Notebook already executed
- [ ] Outputs visible
- [ ] No dependency on reviewer rerunning notebook
- [ ] README included
- [ ] No plagiarism/copy-paste implementation

---

## 31. Priority Roadmap

### Phase 1: Data Collection — P0

```text
Choose source
→ Build scraper
→ Collect 15k–20k raw reviews
→ Save raw data
```

### Phase 2: Dataset Preparation — P0

```text
Clean
→ Deduplicate
→ Label
→ Analyze distribution
→ Produce 10k+ usable records
```

### Phase 3: Baseline — P0

```text
TF-IDF
→ SVM
→ Evaluation
```

### Phase 4: Deep Learning — P0

```text
Tokenizer
→ Embedding
→ BiLSTM
→ Evaluation
```

### Phase 5: Second Deep Model — P1

```text
Embedding
→ CNN
→ BiLSTM
→ Evaluation
```

### Phase 6: Optimization — P1

Tune:

```text
sequence length
embedding dimension
LSTM units
dropout
batch size
learning rate
epochs
```

### Phase 7: Inference — P0

```text
Text
→ Model
→ Probability
→ Sentiment class
```

### Phase 8: Documentation — P0

```text
README
→ requirements
→ executed notebook
→ verify outputs
→ ZIP
```

---

## 32. Definition of Done

```text
✓ Python project
✓ Independent scraping
✓ ≥10,000 clean samples
✓ 3 sentiment classes
✓ Feature extraction
✓ Labeling
✓ 3 training experiments
✓ Deep learning model
✓ Best model >92% target
✓ Minimum testing accuracy ≥85%
✓ Evaluation metrics
✓ Confusion matrix
✓ Inference implementation
✓ Categorical inference output
✓ requirements.txt
✓ scraping code
✓ executed training.ipynb
✓ dataset CSV
✓ README
✓ single ZIP submission
```

---

## 33. Final Target Architecture

```text
               GOOGLE PLAY REVIEWS
                       │
                       ▼
                PYTHON SCRAPER
                       │
                       ▼
              15K–20K RAW REVIEWS
                       │
                       ▼
             CLEAN + DEDUPLICATE
                       │
                       ▼
                RATING LABELING
                       │
             ┌─────────┼─────────┐
             ▼         ▼         ▼
          Negative   Neutral   Positive
             └─────────┼─────────┘
                       ▼
                 10K+ DATASET
                       │
                       ▼
                 80 / 10 / 10
                       │
        ┌──────────────┼──────────────┐
        ▼              ▼              ▼
    TF-IDF/SVM      BiLSTM       CNN+BiLSTM
        │              │              │
        └──────────────┼──────────────┘
                       ▼
                 MODEL COMPARISON
                       │
                       ▼
                 BEST MODEL >92%
                       │
                       ▼
                   INFERENCE
                       │
              ┌────────┼────────┐
              ▼        ▼        ▼
           NEGATIVE  NEUTRAL  POSITIVE
```

---

## 34. Strategic Recommendation

Build around:

**Google Play reviews + rating-derived 3-class labels + TF-IDF/SVM baseline + BiLSTM + CNN/BiLSTM.**

This provides a direct path toward:

- all mandatory criteria,
- 10,000+ samples,
- 3 classes,
- deep learning,
- 3 experiments,
- inference,
- and the Rating 5 recommendation set.

The critical constraint is not writing the neural network. It is obtaining **10,000+ clean, independently scraped, sufficiently diverse Indonesian reviews** and producing a test result that genuinely generalizes.
