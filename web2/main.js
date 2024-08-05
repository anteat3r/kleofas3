document.querySelector("#sub").addEventListener('click', async (_event) => {
  await navigator.serviceWorker.register("service-worker.js");
  await Notification.requestPermission();
  const reg = await navigator.serviceWorker.ready;
  const subscription = await reg.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey: "BIJ29i59x2PSDgMBTTMnYW5lQjStAMrbRAGDmgcgT26iWcRmK5GFjJ1oUAVtL_oiOMwVxEsMjX2z5ASZ_PMziFE",
  })
  console.log(JSON.stringify(subscription));
});
