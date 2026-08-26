# Indonesian Sentiment Analysis — Project Workspace
Locked via PLAN.md (Q1–Q18). See ../PLAN.md and ../PRD.md

## Structure
```
project/
  scraping/scrape.py
  notebooks/training.ipynb  (to be created, 18 sections §19)
  dataset/sentiment_dataset.csv  (15–20k raw → 10k+ clean)
  src/preprocessing.py, labeling.py, evaluation.py, inference.py
  results/confusion_matrix.png, history.png
  requirements.txt
```

## Quick start
```bash
cd project/scraping && python scrape.py
cd ../notebooks && jupyter notebook training.ipynb  # Run All on Colab T4
```

## Provenance
- Apps: MLBB, Gojek, Shopee, Tokopedia
- Labeling: 1–2→negative, 3→neutral, 4–5→positive (weak supervision)
- Split: 80/10/10 stratified, fit on train only
