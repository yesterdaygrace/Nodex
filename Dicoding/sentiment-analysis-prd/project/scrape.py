"""
Indonesian Sentiment Analysis — Google Play Scraper
Locked decisions: Q1 multi-app, Q2 15-20k raw paginated, Q5 decoupled.
Apps: com.mobile.legends, com.gojek.app, com.shopee.id, com.tokopedia.tkpd
"""
from google_play_scraper import reviews
import pandas as pd
import time, random

APPS = [
    ("com.mobile.legends", "MLBB"),
    ("com.gojek.app", "Gojek"),
    ("com.shopee.id", "Shopee"),
    ("com.tokopedia.tkpd", "Tokopedia"),
]
TARGET_PER_APP = 5000
LANG = "id"
COUNTRY = "id"
OUT = "../dataset/sentiment_dataset.csv"
RAW_OUT = "../dataset/raw_reviews.csv"

def scrape_app(app_id, target=5000, lang="id", country="id"):
    all_reviews = []
    token = None
    while len(all_reviews) < target:
        batch, token = reviews(app_id, lang=lang, country=country, count=200, continuation_token=token)
        if not batch:
            break
        all_reviews.extend(batch)
        if token is None:
            break
        time.sleep(random.uniform(0.5, 1.2))
        if len(all_reviews) % 1000 == 0:
            print(f"  {app_id}: {len(all_reviews)}/{target}")
    return all_reviews[:target]

if __name__ == "__main__":
    rows = []
    for app_id, app_name in APPS:
        print(f"Scraping {app_name} ({app_id})...")
        data = scrape_app(app_id, TARGET_PER_APP, LANG, COUNTRY)
        for r in data:
            rows.append({
                "review_id": r.get("reviewId"),
                "text": r.get("content"),
                "rating": r.get("score"),
                "date": r.get("at"),
                "app_id": app_id,
                "app_name": app_name,
            })
        print(f"  -> {len(data)} collected")
    df = pd.DataFrame(rows)
    df = df.dropna(subset=["text"]).drop_duplicates(subset=["text"])
    df.to_csv(RAW_OUT, index=False)
    print(f"Raw saved: {RAW_OUT} ({len(df)} rows)")

    # quick weak labeling for audit
    def label(s):
        if s <= 2: return "negative"
        if s == 3: return "neutral"
        return "positive"
    df["label"] = df["rating"].apply(label)
    df = df[["review_id","text","rating","label","app_id","app_name","date"]]
    df.to_csv(OUT, index=False)
    print(f"Labeled saved: {OUT} ({len(df)} rows)")
    print(df["label"].value_counts())
