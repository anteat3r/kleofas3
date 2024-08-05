/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />
const sw = /** @type {ServiceWorkerGlobalScope} */ (/** @type {unknown} */ (self));

/* eslint-env browser, serviceworker */

sw.addEventListener('install', () => {
  sw.skipWaiting();
})

sw.addEventListener('push', (event) => {
  const datatext = event.data.text();
  const data = JSON.parse(datatext);

  if (data.type === "notif") {
    event.waitUntil(
      sw.registration.showNotification(data.title)
    )
  }

  localStorage.

})
