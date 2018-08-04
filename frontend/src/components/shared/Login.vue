<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><span class="fa fa-lock"></span> Login</h1>
                </div>
            </div>
        </div>
        <section class="section">
            <div class="container">
                <form v-on:submit.prevent="login">
                    <b-field label="Username"
                             :type="errors.has('Username') ? 'is-danger' : ''"
                             :message="errors.first('Username')">
                        <b-input v-model.trim="username"
                                 name="Username"
                                 v-validate="'required'">
                        </b-input>
                    </b-field>

                    <b-field label="Password"
                             :type="errors.has('Password') ? 'is-danger' : ''"
                             :message="errors.first('Password')">
                        <b-input type="password" v-model="password"
                                 name="Password"
                                 v-validate="'required'"
                                 password-reveal>
                        </b-input>
                    </b-field>

                    <button :disabled="errors.any()"
                            v-bind:class="{'is-loading': loading}"
                            class="button is-primary">Login
                    </button>
                </form>
            </div>
        </section>
    </div>
</template>

<script>
    export default {
        data() {
            return {
                username: '',
                password: '',
                loading: false
            };
        },
        methods: {
            login: function () {
                this.$validator.validateAll().then((result) => {
                    if (result) {
                        this.loading = true;

                        this.$http.post(this.$store.state.backendURL + '/user/login', {
                            username: this.username,
                            password: this.password
                        }).then(res => {
                            this.loading = false;

                            // Decode JWT
                            const base64Url = res.body.token.split('.')[1];
                            const base64 = base64Url.replace('-', '+').replace('_', '/');
                            const userData = JSON.parse(window.atob(base64));

                            this.$store.commit('setUser', {
                                user: {
                                    name: userData.id,
                                    token: res.body.token,
                                    exp: userData.exp
                                }
                            });

                            this.$toast.open({
                                type: 'is-success',
                                message: 'The Login was successfull.'
                            });

                            this.$router.push({path: '/'});
                        }, () => {
                            this.loading = false;
                        });
                    }
                });
            }
        }
    };
</script>