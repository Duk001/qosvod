<template>
  <div class="categories-select-body">
    <ul class="list">
      <li
        v-for="(data, index) in listOfCategories"
        :key="`data-${index}`"
        :class="[
          selectedCategories.includes(data)
            ? 'film-category-elem-selected'
            : 'film-category-elem',
        ]"
        v-on:click="toggleCategory(data)"
      >
        {{ data }}
      </li>
    </ul>
    <div class="cat-buttons">
      <div class="reset-button" v-on:click="resetCategory()">Reset</div>
      <div class="close-button" v-on:click="closeCategory()">Close</div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
export default {
  inject: ["host"],
  data() {
    return {
      listOfCategories: [],
      selectedCategories: [],
    };
  },
  async created() {
    // GET request using axios with error handling
    await axios
      .get(this.host + "categories")
      .then((response) => (this.listOfCategories = response.data))
      .catch((error) => {
        this.errorMessage = error.message;
        console.error("There was an error!", error);
      });
    //console.log(this.listOfCategories);
  },

  methods: {
    toggleCategory(category) {
      if (!this.selectedCategories.includes(category)) {
        this.selectedCategories.push(category);
      } else {
        var index = this.selectedCategories.indexOf(category);
        if (index !== -1) {
          this.selectedCategories.splice(index, 1);
        }
      }

      //console.log(this.selectedCategories);
      this.$emit("listOfCategories", this.selectedCategories);
    },
    resetCategory() {
      this.selectedCategories = [];
      this.$emit("listOfCategories", this.selectedCategories);

    },
    closeCategory() {
      this.$emit("changeCategoriesVisibility", false);
      // this.$emit("listOfCategories", this.selectedCategories);
    },
  },
};
</script>

<style>
.categories-select-body {
  border: 3px solid #444;
  border-radius: 5px;
  /* border-bottom: 3px solid #444; */

  position: fixed;
  z-index: 3;
  list-style-type: none;
  margin: 50px 20% 0 20%;
  padding: 10px;
  /* height: 20px; */
  width: 60%;
  overflow: hidden;
  background-color: #222;
  -webkit-user-select: none; /* Chrome all / Safari all */
  -moz-user-select: none; /* Firefox all */
  -ms-user-select: none; /* IE 10+ */
  user-select: none; /* Likely future */
}
.film-category-elem {
  margin: 3px;
  padding: 8px;
  color: aliceblue;
  border: 3px solid #555;
  border-radius: 20px;
  cursor: pointer;
}
.film-category-elem-selected {
  margin: 3px;
  padding: 8px;
  color: aliceblue;
  background: #555;
  border: 3px solid #555;
  border-radius: 20px;
  cursor: pointer;
}
.cat-buttons {
  margin: 5px;
  color: aliceblue;

  display: grid;
  grid-template-columns: 50% 1fr;
}
.close-button {
  margin: auto;
  padding: 8px;
  background: #555;
  border: 3px solid #555;
  border-radius: 20px;
  cursor: pointer;

  width: 50%;
  align-self: center;
  float: right;
}
.reset-button {
  margin: auto;
  padding: 8px;
  background: #555;
  border: 3px solid #555;
  border-radius: 20px;
  cursor: pointer;

  width: 50%;
  float: left;
}
</style>