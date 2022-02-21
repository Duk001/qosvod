<template>
  <div class="body">
    <video-player>
    </video-player>
  </div>
</template>
<script>
//import MainLayout from '../components/Main.vue'
import VideoPlayer from '../components/VideoPlayer.vue'
export default {
  inject:['host'],
  components: {
      VideoPlayer
}
}



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
html,body { margin: 0 !important;
background-color: #000000;}
</style>
