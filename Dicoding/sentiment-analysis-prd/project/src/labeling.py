def rating_to_label(rating: int) -> str:
    if rating <= 2: return "negative"
    if rating == 3: return "neutral"
    return "positive"
