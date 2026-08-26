"""Minimal cleaning — Q6/Q14 locked: keep negation, demojize, no stemming by default."""
import re, emoji

SLANG = {
    "gk":"tidak","ga":"tidak","nggak":"tidak","ngga":"tidak","enggak":"tidak",
    "yg":"yang","bgt":"banget","bngt":"banget","bangettt":"banget",
    "wkwk":"tertawa","wkwkwk":"tertawa","mantul":"mantap betul",
    "cuy":"kawan","bro":"kawan","sis":"kawan","gw":"saya","gue":"saya",
    "lo":"kamu","lu":"kamu","km":"kamu","aq":"aku","aku":"aku",
    "tdk":"tidak","gak":"tidak","ngk":"tidak","jgn":"jangan",
    "udh":"sudah","udah":"sudah","blm":"belum","belum":"belum",
    "krn":"karena","karna":"karena","klo":"kalau","kl":"kalau",
    "skrg":"sekarang","bs":"bisa","bisa":"bisa","lagi":"lagi",
    "aja":"saja","ajaa":"saja","doi":"dia","kpn":"kapan",
}

NEGATION_KEEP = {"tidak","jangan","belum","bukan","tanpa","kurang"}

def clean_minimal(text: str) -> str:
    t = str(text).lower()
    t = re.sub(r"http\S+|www\.\S+", " ", t)
    t = emoji.demojize(t, delimiters=(" ", " "))
    t = re.sub(r"@\w+|#\w+", " ", t)
    t = re.sub(r"(.)\1{2,}", r"\1\1", t)  # tttt -> tt
    t = re.sub(r"\s+", " ", t).strip()
    toks = [SLANG.get(w, w) for w in t.split()]
    return " ".join(toks)
