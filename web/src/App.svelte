<script lang="ts">
  import Home from './lib/Home.svelte'
  import githubLogo from '../public/github.svg'
  import discordLogo from '../public/discord.svg'
  import pb from './pb'

  async function loginGithub() {
    await pb.collection("users")
        .authWithOAuth2({ provider: "github" })
    unique = {}
  }
  async function loginDiscord() {
    await pb.collection("users")
        .authWithOAuth2({
          provider: "discord",
          scopes: ["identify", "email"],
        })
    unique = {}
  }

  let unique = {}
</script>

{#key unique}
  {#if pb.authStore.isValid}
    <Home />
  {:else}
    <button on:click={loginGithub}>
      <img src={githubLogo} alt="github logo">
      Login with Github
    </button>
    <button on:click={loginDiscord}>
      <img src={discordLogo} alt="discord logo">
      Login with Discord
    </button>
  {/if}
{/key}

