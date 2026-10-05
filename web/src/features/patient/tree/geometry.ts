type Point = [number, number]

export interface TwoTone {
  light: string
  shade: string
}

const curveSamples = 18
const lineSamples = 5
const light: Point = [-Math.SQRT1_2, -Math.SQRT1_2]

export function branchCenterline(from: Point, to: Point, bend: number): string {
  const [c1, c2, end] = branchControls(from, to, bend)
  return `M${pt(from)} C${pt(c1)} ${pt(c2)} ${pt(end)} L${pt(to)}`
}

export function branchRibbon(from: Point, to: Point, bend: number, startWidth: number, endWidth: number): TwoTone {
  const [c1, c2, end] = branchControls(from, to, bend)
  const points: Point[] = []
  for (let i = 0; i <= curveSamples; i++) {
    points.push(cubic(from, c1, c2, end, i / curveSamples))
  }
  for (let i = 1; i <= lineSamples; i++) {
    points.push(lerpPoint(end, to, i / lineSamples))
  }

  const lengths = cumulativeLengths(points)
  const total = lengths[lengths.length - 1] || 1
  const direction: Point = [to[0] - from[0], to[1] - from[1]]
  const litSide = -direction[1] * light[0] + direction[0] * light[1] > 0 ? 1 : -1

  const lit: Point[] = []
  const dark: Point[] = []
  points.forEach((point, i) => {
    const [nx, ny] = normalAt(points, i)
    const width = startWidth + (endWidth - startWidth) * (lengths[i] / total)
    lit.push([point[0] + nx * width * litSide, point[1] + ny * width * litSide])
    dark.push([point[0] - nx * width * litSide, point[1] - ny * width * litSide])
  })

  return {
    light: polygon([...lit, ...[...points].reverse()]),
    shade: polygon([...points, ...[...dark].reverse()]),
  }
}

export interface TrunkGeometry {
  x: number
  tipY: number
  rootY: number
  groundY: number
  topWidth: number
  baseWidth: number
  flare: number
}

export function trunkShape({ x, tipY, rootY, groundY, topWidth, baseWidth, flare }: TrunkGeometry): TwoTone {
  const half = (side: number) =>
    `M${pt([x, tipY])} ` +
    `Q${pt([x + side * topWidth, tipY])} ${pt([x + side * topWidth, tipY + topWidth * 2.2])} ` +
    `L${pt([x + side * baseWidth, rootY])} ` +
    `Q${pt([x + side * baseWidth, groundY])} ${pt([x + side * (baseWidth + flare), groundY])} ` +
    `L${pt([x, groundY])} Z`
  return { light: half(-1), shade: half(1) }
}

function branchControls(from: Point, to: Point, bend: number): [Point, Point, Point] {
  const rise = Math.max(from[1] - to[1], 1)
  const turn = Math.min(rise, bend)
  return [
    [from[0], from[1] - turn * 0.5],
    [to[0], from[1] - turn * 0.5],
    [to[0], from[1] - turn],
  ]
}

function cubic(p0: Point, p1: Point, p2: Point, p3: Point, t: number): Point {
  const u = 1 - t
  const a = u * u * u
  const b = 3 * u * u * t
  const c = 3 * u * t * t
  const d = t * t * t
  return [a * p0[0] + b * p1[0] + c * p2[0] + d * p3[0], a * p0[1] + b * p1[1] + c * p2[1] + d * p3[1]]
}

function lerpPoint(a: Point, b: Point, t: number): Point {
  return [a[0] + (b[0] - a[0]) * t, a[1] + (b[1] - a[1]) * t]
}

function cumulativeLengths(points: Point[]): number[] {
  const lengths = [0]
  for (let i = 1; i < points.length; i++) {
    lengths.push(lengths[i - 1] + Math.hypot(points[i][0] - points[i - 1][0], points[i][1] - points[i - 1][1]))
  }
  return lengths
}

function normalAt(points: Point[], i: number): Point {
  const prev = points[Math.max(i - 1, 0)]
  const next = points[Math.min(i + 1, points.length - 1)]
  const dx = next[0] - prev[0]
  const dy = next[1] - prev[1]
  const length = Math.hypot(dx, dy) || 1
  return [-dy / length, dx / length]
}

function polygon(points: Point[]): string {
  return `M${points.map(pt).join(' L')} Z`
}

function pt([x, y]: Point): string {
  return `${x.toFixed(1)} ${y.toFixed(1)}`
}
