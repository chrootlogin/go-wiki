import {GlobalComponents, LocalComponents, WikiComponents} from './components';
import Vue from 'vue'
import VueRouter from 'vue-router';
import store from './store/index';

const routes = [
    /*{
        path: '/',
        component: LocalComponents.Home
    },*/
    {
        path: '/register',
        component: GlobalComponents.Register
    },
    {
        path: '/login',
        component: GlobalComponents.Login
    },
    {
        path: '/profile',
        component: GlobalComponents.Profile
    },
    {
        path: '/createSpace',
        component: WikiComponents.CreateSpace
    },
    /*{
        path: '/wiki',
        component: WikiComponents.ListSpaces
    },*/
    {
        name: 'page',
        path: '/:pageSlug*',
        component: WikiComponents.GetPage,
        props: true
    }
];

const router = new VueRouter({routes});

router.beforeEach((to, from, next) => {
    // Cleanup old notifications
    store.commit('setNotification', {notification: {}});

    // if not logged in log-in as anonymous
    /*if(
        !store.state.user
        || (store.state.user && store.state.user.exp < Date.now() / 1000)
    ) {
        Vue.http.post(store.state.backendURL + "/user/login", {
            username: "anonymous",
            password: "anonymous"
        }).then(res => {
            // Decode JWT
            const base64Url = res.body.token.split('.')[1];
            const base64 = base64Url.replace('-', '+').replace('_', '/');
            const userData = JSON.parse(window.atob(base64));

            store.commit('setUser', {
                user: {
                    name: userData.id,
                    token: res.body.token,
                    exp: userData.exp
                }
            });

            next();
        });

    }*/

    next();
});

export default router;

