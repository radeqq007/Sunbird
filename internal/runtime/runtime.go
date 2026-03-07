package runtime

import _ "embed"

//go:embed shared.ts
var SharedRuntimeSrc string

//go:embed node.ts
var NodeRuntimeSrc string

//go:embed deno.ts
var DenoRuntimeSrc string

//go:embed bun.ts
var BunRuntimeSrc string

//go:embed web.ts
var WebRuntimeSrc string
