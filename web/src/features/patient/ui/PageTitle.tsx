interface PageTitleProps {
  main: string
  rest: string
}

export function PageTitle({ main, rest }: PageTitleProps) {
  return (
    <h1 className="text-2xl font-medium tracking-tight lg:hidden">
      {main} <span className="text-muted">{rest}</span>
    </h1>
  )
}
