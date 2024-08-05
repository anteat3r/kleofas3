<script lang="ts">
  import Home from './lib/Home.svelte'
  import githubLogo from '../public/github.svg'
  import pb from './pb'

  async function login() {
    await pb.collection("users")
        .authWithOAuth2({ provider: "github" })
    unique = {}
  }

  let unique = {}
</script>

{#key unique}
  {#if pb.authStore.isValid}
    <Home />
  {:else}
    <button on:click={login}>
      <img src={githubLogo} alt="github logo">
      Login with Github
    </button>
  {/if}
{/key}

