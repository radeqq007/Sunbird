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

