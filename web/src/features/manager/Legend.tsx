interface LegendItem {
  color: string
  label: string
}

interface LegendProps {
  items: LegendItem[]
}

export function Legend({ items }: LegendProps) {
  return (
    <div className="flex flex-wrap gap-4 text-sm text-subtle">
      {items.map((item) => (
        <span key={item.label} className="flex items-center gap-2">
          <span className="size-3 rounded-sm" style={{ backgroundColor: item.color }} aria-hidden />
          {item.label}
        </span>
      ))}
    </div>
  )
}
