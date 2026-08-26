# Plan — Indonesian Sentiment Analysis (Dicoding Rating 5)
*Grilled 2026-08-25 — 18 decisions, 3 rounds, frontier empty*

## Shared Understanding

**Goal:** Ship `sentiment-analysis.zip` containing `training.ipynb` (executed), `scrape.py`, `sentiment_dataset.csv`, `requirements.txt`, `README.md` that passes Dicoding review at Rating 5: 10k+ clean Indonesian reviews, 3 classes, 3 experiments (TF-IDF/SVM + BiLSTM + CNN+BiLSTM), deep learning, inference visible, ≥85% test (target >92% honest).

**Non-goals:** No Kaggle dataset masquerading as scrape. No leakage-inflated 92%. No training-time network calls.

---

## Design Tree — Decisions Locked

| # | Decision | Locked Choice | Rationale |
|---|----------|---------------|-----------|
| Q1 | Data domain | **Multi-app: MLBB + Gojek + Shopee + Tokopedia** (`com.mobile.legends`, `com.gojek.app`, `com.shopee.id`, `com.tokopedia.tkpd`) | Lexical diversity, 4×5k mitigates single-app bias |
| Q2 | Scale | **15–20k raw → 10k+ clean** via `continuation_token` loop, dedup on `content` | Accounts 15–25% loss; `count=5000` alone caps |
| Q3 | Labeling | **Rating-derived weak supervision** 1-2 Neg, 3 Neu, 4-5 Pos + **200-sample manual audit** | Honest about noise (~6% error); reviewer expects disclaimer |
| Q4 | Model lever for >92% | **Keep 3-experiment shape, upgrade embeddings:** E2 scratch 128d, E3 FastText-id 300d frozen + `class_weight=balanced` | +2–3pp without breaking PRD diagram |
| Q5 | Reproducibility contract | **Decoupled:** `scrape.py` writes CSV; `training.ipynb` only `read_csv`, no network | `Run All` deterministic; date provenance in README |
| Q6 | Preprocessing | **Minimal `clean_minimal()`** — keep negation, keep digits/punct, demojize emojis, 60-entry slang map, **no Sastrawi/stopword by default**; ablation only | Preserves sentiment signals, Neutral especially |
| Q7 | Hardware/stack | **Colab T4 GPU**, `tensorflow==2.15.0`, `scikit-learn`, `pandas`, `numpy<2`, `Sastrawi`, `emoji`, `wordcloud`, `google-play-scraper`; `pipreqs` + prune | Fits 3 experiments in <45 min |
| Q8 | Imbalance | **Stratified 80/10/10 + `class_weight=balanced`**, no SMOTE | Neutral ~15%; F1 per class proves understanding |
| Q9 | Tokenizer | **`Tokenizer(20000, oov="<OOV>")`, `max_len=128`, post pad/trunc**, emb 128 vs 300 | Median review 18 tok; covers OOV <4% |
| Q10 | Split/leakage | **Single 80/10/10 stratified for all**; vectorizer/tokenizer `fit` on train only | PRD §17 pipeline verbatim |
| Q11 | Evaluation | **Per-experiment: train/val/test acc + macro/per-class P/R/F1 + `confusion_matrix.png` + `history.png`**, `EarlyStopping(patience=3)` | Rating 5 evidence |
| Q12 | Inference | **`def predict_sentiment(text)->str`** + 6 demo calls (2/class) with probs visible in notebook | Categorical + probs |
| Q13 | ZIP | **Flat ZIP root:** `training.ipynb`, `scrape.py`, `sentiment_dataset.csv`, `requirements.txt`, `README.md`; notebook loads `sentiment_dataset.csv` flat (fallback `dataset/...`) | Matches Dicoding uploader expectation |
| Q14 | Clean rules | Lower → URL strip → demojize → repeat-char limit 2 → `@/#` strip → slang map → whitespace norm | Locked `clean_minimal()` |
| Q15 | <92% fallback | **Ship honest 90–91% if that's true**; appendix "Path to 92%+ (IndoBERT fine-tune)" | Never leak; 85% still passes |
| Q16 | Repro cell | **Seed cell top of notebook:** `PYTHONHASHSEED+random+np+tf=42`, version prints, provenance block | Rating 5 audit |
| Q17 | Build order | **Waterfall:** Scrape → Clean/Audit → E1 → E2 → E3 → Eval → Inference → ZIP | E1 informs E2/E3 hyperparams |
| Q18 | Artifact | **This PLAN.md + 8 tasks** | Enables `/gsd-plan-phase` |

---

## Architecture (final target)

```
Google Play reviews (4 apps, paginated)
  → 15–20k raw (review_id, content, score, at, app_id)
  → clean_minimal + dedup + empty removal
  → rating labeling + 200 audit
  → 10k+ dataset (id, text, rating, label, source, created_at)
  → stratified 80/10/10
  → ┌─ E1: TF-IDF(15k, 1-2gram) + LinearSVC(GridSearch C, 5-fold)
  │  E2: Tokenizer(20k)→Pad128→Emb128→BiLSTM128→Drop0.5→Dense3
  │  E3: Tokenizer(20k)→Pad128→Emb300 FastText frozen→Conv1D128→BiLSTM128→Dense3
  → Eval table + confusion + history
  → best model → predict_sentiment()
  → executed training.ipynb + flat ZIP
```

## Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| >92% not hit on honest split | High | Pretrained emb + class_weight; ship honest 91% + appendix |
| `google-play-scraper` throttling | Med | Backoff + continuation_token, run once, cache CSV |
| Neutral F1 collapse | Med | Keep negation, no stemming, class_weight, report per-class |
| Clean over-aggressive | Med | Minimal pipeline, ablation appendix |
| Notebook not `Run All` clean | High | No network in training, flat CSV path, executed before ZIP |

## 8-Phase Roadmap (P0/P1)

| Phase | Name | Done when |
|-------|------|-----------|
| 1 | Data Collection — P0 | `scrape.py` paginates 4 apps → `raw_reviews.csv` 15k+ |
| 2 | Dataset Prep — P0 | `dataset/sentiment_dataset.csv` 10k+ clean + audit table + distribution plot |
| 3 | Baseline E1 — P0 | TF-IDF+SVM GridSearch, 80/10/10, F1+CM |
| 4 | Deep E2 — P0 | BiLSTM scratch, EarlyStopping, history.png |
| 5 | Deep E3 — P1 | CNN+BiLSTM FastText, comparison table |
| 6 | Optimization — P1 | Tune max_len/emb/dropout/batch if <92% honest |
| 7 | Inference — P0 | `predict_sentiment` + 6 demos visible |
| 8 | Documentation — P0 | README 16 sections, requirements.txt, executed notebook, flat ZIP verify `Run All` |

## File Contract

```
sentiment-analysis-prd/
  PLAN.md              ← this file
  PRD.md / README.md
  scrape.py            ← paginated 4-app scraper
  training.ipynb       ← executed, 18-section structure (§19)
  sentiment_dataset.csv← 10k+ flat (also dataset/sentiment_dataset.csv)
  requirements.txt     ← pipreqs pruned
  README.md            ← provenance + 16 sections (§28)
  results/ (optional)  ← confusion_matrix.png, history.png
```

## Definition of Done (§32)

All 21 checks pass including 3 experiments, 2 dimensions changed, deep learning, ≥85% test (target >92% honest), per-class metrics, confusion matrix, inference 3 classes visible, seeds, flat ZIP, no plagiarism, indie scrape documented.

## Next Action

Confirm this PLAN, then `/gsd-plan-phase 1` to scaffold scraper. Or say "tweak Qx" to reopen a branch.
