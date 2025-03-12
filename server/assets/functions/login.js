const login_page = {
  // Variables Section
  login: true,
  error: false,
  forgot: false,
  otp: false,

  login_data: [
    {
      name: "username",
      label: "Username",
      complete: "username",
      type: "text",
      placeholder: "eg. user_name_123",
    },
    {
      name: "password",
      label: "Password",
      complete: "current-password",
      type: "password",
      placeholder: "eg. xxxxxxxx",
    },
  ],

  signup_data: [
    {
      name: "username",
      label: "Username",
      complete: "username",
      type: "text",
      placeholder: "eg. user_name_123",
    },
    {
      name: "email",
      label: "Email",
      complete: "email",
      type: "text",
      placeholder: "eg. user@email.com",
    },
    {
      name: "password",
      label: "Password",
      complete: "new-password",
      type: "password",
      placeholder: "eg. xxxxxxxx",
    },
    {
      name: "confirm-password",
      label: "Confirm Password",
      complete: "new-password",
      type: "password",
      placeholder: "eg. xxxxxxxx",
    },
  ],

  // Functions Section
  title() {
    sessionStorage.setItem("login", this.login);
    document.title = `${this.login ? "Login" : "Sign Up"} | ${document.title}`;
  },

  get_state() {
    const data = sessionStorage.getItem("login");
    if (!data) return;
    this.login = data === "true";
  },

  error_message(error) {
    this.error = true;
    this.$refs.error_msg.textContent = error.message;
    setTimeout(() => {
      this.error = false;
    }, 5000);
  },

  async create_form(el) {
    try {
      let formData = new FormData(el);
      const password = formData.get("password");
      const c_password = formData.get("confirm-password");
      if (password) formData.set("password", await this.hash_data(password));
      if (c_password)
        formData.set("confirm-password", await this.hash_data(c_password));
      return formData;
    } catch (error) {
      this.error_message(error);
    }
  },

  async hash_data(value) {
    const encoder = new TextEncoder();
    const data = encoder.encode(value);
    const hash_buff = await crypto.subtle.digest("SHA-256", data);
    return Array.from(new Uint8Array(hash_buff))
      .map((byte) => byte.toString(16).padStart(2, "0"))
      .join("");
  },

  async form_submit() {
    try {
      const formData = await this.create_form(this.$el);
      let response = await fetch(this.login ? "/login" : "/signup", {
        method: "POST",
        body: formData,
      });
      if (!response.ok) throw new Error(await response.text());
      const auth = response.headers.get("X-Auth-Token");
      if (!auth) throw new Error("Auth token not found!");
      const { token, route } = JSON.parse(auth);
      if (!token || !route)
        throw new Error("Login failed with no data or route");
      sessionStorage.setItem("token", token);
      sessionStorage.setItem("login", this.login ? "true" : "false");
      window.location.href = route;
    } catch (error) {
      this.error_message(error);
    }
  },

  async google_login() {
    try {
      let res = await fetch("/google", {
        method: "POST",
      });
      if (!res.ok) throw new Error(await res.text());
      const data = await res.json();
      window.location.href = data.url;
    } catch (error) {
      this.error_message(error);
    }
  },
};
