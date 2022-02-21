<template>
<div>
  <!-- <nav-bar>
  </nav-bar> -->
  <h1 class="msg">Page not found.   :(</h1>
</div>
</template>

<script>
import NavBar from '../components/NavBar.vue'
export default {
  inject:['host'],
  components: {
    NavBar
  }
}
</script>
<style scoped>
.msg{
    color: rgb(255, 255, 255);
}
</style>