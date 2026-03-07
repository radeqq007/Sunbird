import * as nodeFs from "fs";
import * as nodePath from "path";

export const fs = {
  read(path: string): string {
    try {
      return nodeFs.readFileSync(path, "utf-8");
    } catch (e) {
      throw new Error(`read: ${(e as Error).message}`);
    }
  },

  write(path: string, content: string): void {
    try {
      nodeFs.writeFileSync(path, content, "utf-8");
    } catch (e) {
      throw new Error(`write: ${(e as Error).message}`);
    }
  },

  append(path: string, content: string): void {
    try {
      nodeFs.appendFileSync(path, content, "utf-8");
    } catch (e) {
      throw new Error(`append: ${(e as Error).message}`);
    }
  },

  remove(path: string): void {
    try {
      if (nodeFs.existsSync(path) && nodeFs.statSync(path).isDirectory()) {
        nodeFs.rmSync(path, { recursive: true, force: true });
      } else {
        nodeFs.unlinkSync(path);
      }
    } catch (e) {
      throw new Error(`remove: ${(e as Error).message}`);
    }
  },

  exists(path: string): boolean {
    try {
      return nodeFs.existsSync(path);
    } catch (e) {
      throw new Error(`exists: ${(e as Error).message}`);
    }
  },

  is_dir(path: string): boolean {
    try {
      return nodeFs.statSync(path).isDirectory();
    } catch (e) {
      throw new Error(`is_dir: ${(e as Error).message}`);
    }
  },

  list_dir(path: string): string[] {
    try {
      return nodeFs.readdirSync(path);
    } catch (e) {
      throw new Error(`list_dir: ${(e as Error).message}`);
    }
  },

  create_dir(path: string): void {
    try {
      nodeFs.mkdirSync(path, { recursive: true });
    } catch (e) {
      throw new Error(`create_dir: ${(e as Error).message}`);
    }
  },

  rename(oldPath: string, newPath: string): void {
    try {
      nodeFs.renameSync(oldPath, newPath);
    } catch (e) {
      throw new Error(`rename: ${(e as Error).message}`);
    }
  },

  copy(src: string, dest: string): void {
    try {
      nodeFs.copyFileSync(src, nodePath.resolve(dest));
    } catch (e) {
      throw new Error(`copy: ${(e as Error).message}`);
    }
  },
};

export const io = {
  print(...args: any[]): void {
    process.stdout.write(
      args
        .map((arg) =>
          typeof arg === "object" ? JSON.stringify(arg) : String(arg),
        )
        .join(" "),
    );
  },
  println(...args: any[]): void {
    this.print(...args);
    process.stdout.write("\n");
  },
  read(prompt?: string): string {
    if (prompt) process.stdout.write(prompt);

    const buffer = Buffer.alloc(1024);
    const bytesRead = nodeFs.readSync(0, buffer, 0, 1024);
    return buffer.toString("utf8", 0, bytesRead);
  },
  readln(prompt?: string): string {
    return this.read(prompt).replace(/\r?\n$/, "");
  },
  printf(template: string, ...args: any[]): void {
    process.stdout.write(format(template, args));
  },
  printfn(template: string, ...args: any[]): void {
    process.stdout.write(format(template, args) + "\n");
  },
  sprintf(template: string, ...args: any[]): string {
    return format(template, args);
  },
  clear(): void {
    console.clear();
  },
  beep(): void {
    process.stdout.write("\u0007");
  },
  get args(): string[] {
    return process.argv.slice(2);
  },
};

const format = (str: string, args: any[]): string => {
  let i = 0;
  return str.replace(/{}/g, () => (i < args.length ? String(args[i++]) : "{}"));
};
