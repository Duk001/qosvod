<template>
  <!-- <div> -->
  <div class="navbar">
    <ul>
      <li><a href="/">Home</a></li>
      <li>
        <div class="cat-button" v-on:click="openCategorySelect()">
          Categories
        </div>
      </li>
      <li>
        <div class="add-button" v-if="username=='admin'" v-on:click="toogleNewFilmFrom">Add</div>
      </li>
      <li class="settings" v-on:click="toggleOvelray">
        <div class="bar1"></div>
        <div class="bar2"></div>
        <div class="bar3"></div>
        <!-- <slide-menu></slide-menu> -->
      </li>
    </ul>
  </div>
  <slide-menu :isOverlay="overlay"></slide-menu>
  <new-film-form
    :isVisible="newFilmForm"
    @changeFormVisibility="changeFV($event)"
  ></new-film-form>
  <!-- </div> -->
</template>
<script>
import NewFilmForm from "./NewFilmForm.vue";
import SlideMenu from "./SlideMenu";

export default {
  name: "NavBar",
  inject: ["host"],
  components: {
    SlideMenu,
    NewFilmForm,
  },
  data() {
    return {
      overlay: false,
      newFilmForm: false,
      username: null,
      // categoryOpen : false,
    };
  },
  mounted() {
    this.username = sessionStorage.username;
  },
  methods: {
    toggleOvelray() {
      console.log("Overlay");
      this.overlay = !this.overlay;
      if (this.overlay == true) {
        document.documentElement.style.overflow = "hidden";
      } else {
        document.documentElement.style.overflow = "auto";
      }
    },
    toogleNewFilmFrom() {
      //console.log("Overlay")
      this.newFilmForm = !this.newFilmForm;
      if (this.newFilmForm == true) {
        document.documentElement.style.overflow = "hidden";
      } else {
        document.documentElement.style.overflow = "auto";
      }
    },
    changeFV(visibility) {
      this.newFilmForm = visibility;
      if (this.newFilmForm == true) {
        document.documentElement.style.overflow = "hidden";
      } else {
        document.documentElement.style.overflow = "auto";
      }
    },
    openCategorySelect() {
      // this.categoryOpen = !this.categoryOpen
      this.$emit("changeCategoriesVisibility", true);
      // this.$emit("listOfCategories", this.selectedCategories);
    },
  },
};
</script>
<style scoped >
.navbar {
  border-bottom: 3px solid #444;
  position: fixed;
  z-index: 3;
  list-style-type: none;
  margin: 0;
  padding: 0;
  height: 50px;
  width: 100%;
  overflow: hidden;
  background-color: #333;
}

li {
  list-style-type: none;
  float: left;
}

li a,
.add-button {
  display: block;
  color: white;
  text-align: center;
  padding: 0px 12px;
  text-decoration: none;
  font-weight: bolder;
  transition: 0.4s;
  cursor: pointer;
}
.cat-button {
  display: block;
  color: white;
  text-align: center;
  padding: 0px 12px;
  text-decoration: none;
  font-weight: bolder;
  transition: 0.4s;
  cursor: pointer;
}
.settings {
  display: block;
  cursor: pointer;
  float: right;
  padding: 0px 12px;
}
.bar1,
.bar2,
.bar3 {
  width: 25px;
  height: 3px;
  background-color: rgb(255, 255, 255);
  margin: 4px 0px;
  transition: 0.4s;
}
:where(li a, .add-button, .cat-button):hover:not(.active) {
  font-size: x-large;
  transition: 0.4s;
}

.active {
  background-color: #04aa6d;
}
</style>
