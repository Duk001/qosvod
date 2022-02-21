<template>
  <div class="login-container">
    <div class="login-form-container">
      <form class="login-form">
        <ul class="login-form-list">
          <li>
            <label for="Login">Login:</label>
            <input type="text" v-model="login" id="Login" name="Login" />
          </li>
          <li>
            <label for="Password">Password:</label>
            <input
              type="password"
              v-model="password"
              id="Password"
              name="Password"
            />
            <div class="login-button" v-on:click="getToken()">Login</div>
          </li>
        </ul>
      </form>
    </div>
  </div>
</template>

<script>
import axios from "axios";
export default {
  inject: ["host"],
  data() {
    return {
      login: null,
      password: null,
      token: null,
    };
  },
  methods: {
    async getToken() {
      let tmpToken;
      var Data = `{"login":"${this.login}",
            "password":"${this.password}" }`;
      var jsonData = JSON.parse(Data);
      const res = await axios.post(this.host + "login", jsonData, {
        headers: {
          "Content-Type": "application/json",
        },
      });

      tmpToken = res.data;
      if (tmpToken && res.status == "200") {
        this.token = tmpToken;
        sessionStorage.token = this.token;
        sessionStorage.username = this.login;
        this.$emit("loginToken", this.token);
        console.log("token: ", this.token);
      }
    },
  },
};
</script>

<style>
.login-container {
  position: absolute;
  display: block;
  width: 100%;
  height: 100%;
  background-color: #222;
  z-index: 4;
}
.login-form-container {
  position: fixed;
  top: 30%;
  left: 50%;
  z-index: 2;
  transform: translate(-50%, -50%);
  background-color: #333;
  color: aliceblue;
  border: 3px solid #444;
  border-radius: 5px;
}
.login-form-list {
  align-content: center;
  list-style-type: none;
  padding: 0px 20px;
}
.login-form-list label,
input {
  display: flex;
  flex-direction: column;
  margin-bottom: 15px;
}
.login-form-list label {
  font-size: large;
  font-weight: bolder;
}
.login-button {
  display: inline;
  background-color: #555;
  color: aliceblue;
  border-radius: 5px;
  font-size: larger;
  padding: 3px 20px;
  border: 3px solid #444;
  cursor: pointer;
}
</style>