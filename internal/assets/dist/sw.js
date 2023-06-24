importScripts("https://cdn.jsdelivr.net/gh/golang/go@go1.20.5/misc/wasm/wasm_exec.js");

// for tinygo
// importScripts("/dist/wasm_exec.js");

function registerWasmHTTPListener(wasm, { urlMatch = (url) => true, args = [] } = {}) {
  let path = new URL(registration.scope).pathname;

  const handlerPromise = new Promise(setHandler => {
    self.wasmhttp = {
      path,
      setHandler,
    };
  });

  const go = new Go();
  go.argv = [ wasm, ...args ];
  WebAssembly.instantiateStreaming(fetch(wasm), go.importObject)
    .then(({ instance }) => go.run(instance));

  addEventListener("fetch", e => {
    const url = new URL(e.request.url);
    if (!urlMatch(url)) return;

    e.respondWith(handlerPromise.then(handler => handler(e.request)));
  });
}

addEventListener("install", (event) => {
  event.waitUntil(skipWaiting());
});

addEventListener("activate", event => {
  event.waitUntil(clients.claim());
});

registerWasmHTTPListener("/dist/client.wasm", {
  urlMatch: (url) => {
    return url.host === "localhost:3000" && !url.pathname.startsWith("/dist");
  }
});
