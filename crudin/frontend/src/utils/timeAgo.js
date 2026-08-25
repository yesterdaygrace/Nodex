// Design §14 relative timestamps: "2 minutes ago", "1 hour ago", "Yesterday", "2 days ago"
export function timeAgo(input) {
  const d = new Date(input)
  if (Number.isNaN(d.getTime())) return "—"

  const now = new Date()
  const sec = Math.floor((now - d) / 1000)
  const min = Math.floor(sec / 60)
  const hr = Math.floor(min / 60)
  const day = Math.floor(hr / 24)

  if (sec < 60) return "just now"
  if (min < 60) return `${min} minute${min === 1 ? "" : "s"} ago`
  if (hr < 24) return `${hr} hour${hr === 1 ? "" : "s"} ago`
  if (day === 1) return "Yesterday"
  if (day < 7) return `${day} day${day === 1 ? "" : "s"} ago`
  return d.toLocaleDateString(undefined, { month: "short", day: "numeric", year: "numeric" })
}
