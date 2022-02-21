<template>
  <div class="container" v-if="isOverlay">
    <div class="overlay" v-on:click="NOisOverlay = false">
      <div class="menu">
        <ul class="menu-list">
          <li>{{username}}</li>
          <!-- <li>Add new film</li> -->
          <li>About</li>
          <li><div class="logout-btn" v-on:click="logout()">Logout</div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
export default {
  inject:['host'],
  props: ["isOverlay"],
  data() {
    return {
      overlay: true,
      newFilmForm: false,
      username : null,
    };
  },
  mounted() {
    this.username = sessionStorage.username
  },
  methods:{
    async logout(){
      let tokenHeader = { token: sessionStorage.token };
      const res = await axios.get(this.host + "logout", {
        headers: tokenHeader,
      });
      console.log(res);
      // if (res.data == sessionStorage.token) {
      //   this.token = sessionStorage.token;
      // } else {
      //   this.token = null;
      // }

      sessionStorage.username = null
      sessionStorage.token = null
      window.location.reload(true)

    }
  }
};
</script>

<style>
.overlay {
  position: absolute;
  display: block;
  width: 100%;
  height: 100%;
  background-color: rgba(51, 51, 51, 0.705);
  z-index: 1;
  cursor: pointer;
}

.menu {
  padding-top: 50px;
  background-color: #333;
  color: aliceblue;
  height: 100%;
  width: 400px;
  transition: 0.4s;
  float: right;
  cursor: auto;
  border: 3px solid #444;
}
.menu-list {
  list-style: none;
}
.logout-btn{
  /* margin: 10px; */
  cursor: pointer;

}
.menu-list li {
  background-color: #333;
  margin: 0px 60px 20px 20px;
  padding: 10px;
  border: 3px solid #444;
  border-radius: 5px;
}
</style>