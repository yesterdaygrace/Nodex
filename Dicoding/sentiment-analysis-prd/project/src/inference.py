"""Inference contract Q12."""
import re, emoji
from src.preprocessing import clean_minimal

def predict_sentiment(text: str, tokenizer, model, max_len=128, labels=("negative","neutral","positive")):
    import numpy as np
    from tensorflow.keras.preprocessing.sequence import pad_sequences
    clean = clean_minimal(text)
    seq = tokenizer.texts_to_sequences([clean])
    pad = pad_sequences(seq, maxlen=max_len, padding="post", truncating="post")
    probs = model.predict(pad, verbose=0)[0]
    idx = int(np.argmax(probs))
    return labels[idx], dict(zip(labels, probs.tolist()))
