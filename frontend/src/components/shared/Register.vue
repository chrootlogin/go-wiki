<template>
    <div>
        <div class="hero is-light">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title"><span class="fa fa-lock"></span> Register</h1>
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
                                 v-validate="'required'">
                        </b-input>
                    </b-field>

                    <b-field label="Email"
                             :type="errors.has('Email') ? 'is-danger' : ''"
                             :message="errors.first('Email')">
                        <b-input type="email" v-model="email"
                                 name="Email"
                                 v-validate="'required'">
                        </b-input>
                    </b-field>

                    <button :disabled="errors.any()"
                            v-bind:class="{'is-loading': loading}"
                            class="button is-primary">Register new account
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
                email: '',
                loading: false
            };
        },
        methods: {
            login: function () {
                this.$validator.validateAll().then((result) => {
                    if (result) {
                        this.loading = true;

                        this.$http.post(this.$store.state.backendURL + '/user/register', {
                            username: this.username,
                            password: this.password,
                            email: this.email
                        }).then(() => {
                            this.loading = false;

                            this.$toast.open({
                                type: 'is-success',
                                message: 'Registration was successfull.'
                            });

                            this.$router.push({path: '/login'});
                        }, () => {
                            this.loading = false;
                        });
                    }
                });
            }
        }
    };
</script>