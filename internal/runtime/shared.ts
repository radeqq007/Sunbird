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
  clamp: (val: number, min: number, max: number): number =>
    Math.min(Math.max(val, min), max),
  sin: (x: number): number => Math.sin(x),
  cos: (x: number): number => Math.cos(x),
  tan: (x: number): number => Math.tan(x),
  pi: Math.PI,
  e: Math.E,
} as const;

export const array = {
  push(arr: unknown[], value: unknown): null {
    arr.push(value);
    return null;
  },

  pop(arr: unknown[]): unknown {
    if (arr.length === 0) throw new Error("RuntimeError: array is empty");
    return arr.pop()!;
  },

  shift(arr: unknown[]): unknown {
    if (arr.length === 0) throw new Error("RuntimeError: array is empty");
    return arr.shift()!;
  },

  unshift(arr: unknown[], value: unknown): null {
    arr.unshift(value);
    return null;
  },

  reverse(arr: unknown[]): null {
    arr.reverse();
    return null;
  },

  index_of(arr: unknown[], value: unknown): number {
    const inspect = (v: unknown): string => JSON.stringify(v) ?? String(v);
    const needle = inspect(value);
    for (let i = 0; i < arr.length; i++) {
      if (inspect(arr[i]) === needle) return i;
    }
    return -1;
  },

  slice(arr: unknown[], start: number, end?: number): unknown[] {
    const e = end ?? arr.length;
    if (start < 0) throw new Error("IndexOutOfBoundsError: Array");
    if (e < 0) throw new Error("IndexOutOfBoundsError: Array");
    if (start > e) throw new Error("RuntimeError: start index is greater than end index");
    return arr.slice(start, e);
  },

  join(arr: unknown[], separator: string): string {
    return arr.map((v) => (typeof v === "string" ? v : JSON.stringify(v))).join(separator);
  },

  concat(arr1: unknown[], arr2: unknown[]): unknown[] {
    return [...arr1, ...arr2];
  },

  contains(arr: unknown[], value: unknown): boolean {
    const inspect = (v: unknown): string => JSON.stringify(v) ?? String(v);
    const needle = inspect(value);
    return arr.some((v) => inspect(v) === needle);
  },

  clear(arr: unknown[]): null {
    arr.splice(0, arr.length);
    return null;
  },
} as const;

// Range helper
export function __range(
  start: number,
  end: number,
  step: number = 1,
): number[] {
  const result: number[] = [];
  if (step > 0) for (let i = start; i < end; i += step) result.push(i);
  else for (let i = start; i > end; i += step) result.push(i);
  return result;
}

// Builtins
export function __len(x: string | unknown[]): number {
  return x.length;
}

export function __append(arr: unknown[], ...items: unknown[]): unknown[] {
  return [...arr, ...items];
}

export function __type(x: unknown): string {
  if (x === null) return "Null";
  if (Array.isArray(x)) return "Array";

  switch (typeof x) {
    case "number":
      return Number.isInteger(x) ? "Integer" : "Float";
    case "string":
      return "String";
    case "boolean":
      return "Boolean";
    case "function":
      return "Function";
    case "object":
      return "Hash";
    default:
      return "Unknown";
  }
}

export function __string(x: unknown): string {
  if (x === null) return "null";

  switch (typeof x) {
    case "string":
      return x;
    case "boolean":
    case "number":
      return x.toString();
    case "object":
      return JSON.stringify(x);
    default:
      return String(x);
  }
}

export function __int(x: unknown): number {
  if (x === null) return 0;

  switch (typeof x) {
    case "number":
      return Math.trunc(x);
    case "string": {
      const n = parseInt(x, 10);
      if (isNaN(n))
        throw new Error(`TypeError: failed to convert string to int: ${x}`);
      return n;
    }
    case "boolean":
      return x ? 1 : 0;
    default:
      throw new Error(
        `TypeError: argument to 'int' not supported, got ${__type(x)}`,
      );
  }
}

export function __float(x: unknown): number {
  if (x === null) return 0.0;

  switch (typeof x) {
    case "number":
      return x;
    case "string": {
      const n = parseFloat(x);
      if (isNaN(n))
        throw new Error(`TypeError: failed to convert string to float: ${x}`);
      return n;
    }
    case "boolean":
      return x ? 1.0 : 0.0;
    default:
      throw new Error(
        `TypeError: argument to 'float' not supported, got ${__type(x)}`,
      );
  }
}

export function __bool(x: unknown): boolean {
  if (x === null) return false;

  switch (typeof x) {
    case "boolean":
      return x;
    case "number":
      return x !== 0;
    case "string":
      return x !== "";
    default:
      return true;
  }
}

export function __exit(code: number = 0): never {
  if (typeof process !== "undefined" && process.exit) {
    process.exit(code);
  }
  // For browsers/other environments
  throw new Error(`exit(${code})`);
}

export function __error(message: string): never {
  throw new Error(`RuntimeError: ${message}`);
}

export function __idx(item: unknown[] | string, index: number): unknown {
  if (index < 0) index = item.length + index;
  if (index < 0 || index >= item.length)
    throw new Error("IndexOutOfBoundsError: Array");
  return item[index];
}

