<template>
  <div class="body">
    <nav-bar> </nav-bar>
    <film-blocks> </film-blocks>
  </div>
</template>

<script>
import FilmBlock from "./components/FilmBlock.vue";
import FilmBlocks from "./components/FilmBlocks.vue";
import NavBar from "./components/NavBar.vue";
export default {
  name: "App",
  components: {
    NavBar,
    FilmBlocks,
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  background-color: #000000;
  margin-top: 0px;
}
html,
body {
  margin: 0 !important;
  background-color: #000000;
}
</style>

