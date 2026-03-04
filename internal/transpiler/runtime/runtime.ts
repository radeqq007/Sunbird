export const io = {
  println: (...args: any) => { console.log(...args) }
} as const 

export const math = {
  abs: (x: number): number => Math.abs(x),
  sqrt: (x: number): number => Math.sqrt(x),
  max: (x: number, y: number): number => Math.max(x, y),
  min: (x: number, y: number): number => Math.min(x, y),
  pow: (x: number, y: number): number => Math.pow(x, y),
  floor: (x: number): number => Math.floor(x),
  ceil: (x: number): number => Math.ceil(x),
  round: (x: number): number => Math.round(x),
  sign: (x: number): number => Math.sign(x),
  clamp: (val: number, min: number, max: number): number => Math.min(Math.max(val, min), max),
  sin: (x: number): number => Math.sin(x),
  cos: (x: number): number => Math.cos(x),
  tan: (x: number): number => Math.tan(x),
  pi: Math.PI,
  e: Math.E,
} as const
