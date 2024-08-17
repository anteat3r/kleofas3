<script lang="ts">
  import pb from '../pb'

  async function setupPush() {
    await navigator.serviceWorker.register("service-worker.js");
    await Notification.requestPermission();
    const reg = await navigator.serviceWorker.ready;
    const subscription = await reg.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: "BIJ29i59x2PSDgMBTTMnYW5lQjStAMrbRAGDmgcgT26iWcRmK5GFjJ1oUAVtL_oiOMwVxEsMjX2z5ASZ_PMziFE",
    })
    console.log(JSON.stringify(subscription));
  }

  let username = ""
  let password = ""
  async function tryLogin() {
    await pb.send("test", {
      method: "POST",
      body: JSON.stringify({
        username,
        password,
      }),
    })
  }

</script>

<button on:click={setupPush}>
  Notif perm request
</button>
<br>
<input type="text" bind:value={username} placeholder="username">
<input type="text" bind:value={password} placeholder="password">
<button on:click={tryLogin}>
  Try Login
</button>
