<template>
  <div class="body">
    <login v-if="this.token == null" @loginToken="assignToken($event)"></login>
    <nav-bar
      v-if="this.token != null"
      @changeCategoriesVisibility="changeCategoriesVisibility($event)"
    >
    </nav-bar>
    <category-select
      v-if="this.token != null && categorySelectVisible"
      @listOfCategories="changeCategories($event)"
      @changeCategoriesVisibility="changeCategoriesVisibility($event)"
    >
    </category-select>
    <film-blocks
      v-if="this.token != null"
      :categories="this.filmCategories"
      :key="filmBlocksKey"
    >
    </film-blocks>
  </div>
</template>

<script>
import axios from "axios";
import FilmBlock from "../components/FilmBlock.vue";
import FilmBlocks from "../components/FilmBlocks.vue";
import NavBar from "../components/NavBar.vue";
import Login from "../components/Login.vue";
import CategorySelect from "../components/CategorySelect.vue";
export default {
  name: "App",
  inject: ["host"],
  components: {
    NavBar,
    FilmBlocks,
    CategorySelect,
    Login,
  },
  data() {
    return {
      token: null,
      login: null,
      categorySelectVisible: false,
      filmCategories: [],
      filmBlocksKey: 0,
    };
    // CategorySelect{
    // }
  },
  async mounted() {
    console.log("Host ", this.host);
    if (sessionStorage.token) {
      //Check if token is valid,
      let tokenHeader = { token: sessionStorage.token };
      const res = await axios.get(this.host + "tokenCheck", {
        headers: tokenHeader,
      });
      console.log(res);
      if (res.data == sessionStorage.token) {
        this.token = sessionStorage.token;
      } else {
        this.token = null;
      }
    }
  },
  methods: {
    assignToken(token) {
      this.token = token;
    },
    changeCategories(categories) {
      this.filmCategories = categories;
      this.filmBlocksKey += 1;
    },
    changeCategoriesVisibility(state) {
      this.categorySelectVisible = state;
    },
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
