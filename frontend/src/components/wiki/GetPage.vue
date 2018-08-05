<template >
    <div v-if="error === 0">
        <article>
            <div class="notification">
                <div class="container has-text-right">
                    <nav class="breadcrump is-pulled-left is-hidden-mobile" aria-label="breadcrumbs">
                        <ul v-html="breadcrumb"></ul>
                    </nav>
                    <button class="button is-success">
                        <span>Create page</span>
                        <span class="icon is-small">
                            <i class="fa fa-plus"></i>
                        </span>
                    </button>
                    <router-link :to="{ name: 'edit', query: { path: url } }" class="button is-primary">
                        <span>Edit page</span>
                        <span class="icon is-small">
                            <i class="fa fa-edit"></i>
                        </span>
                    </router-link>
                    <button v-on:click="deletePage" class="button is-danger is-outlined">
                        <span>Delete page</span>
                        <span class="icon is-small">
                            <i class="fa fa-times"></i>
                        </span>
                    </button>
                </div>
            </div>
            <div class="container" v-if="!edit">
                <div class="content" v-html="page.content"></div>
            </div>
        </article>
    </div>
</template>

<script>
    export default {
        data() {
            return {
                loading: false,
                error: 0,
                page: {
                    title: "",
                    content: ""
                },
                breadcrumb: ""
            }
        },
        computed: {
            url() {
                let url = this.pageSlug;
                if(url == null) {
                    url = "";
                }

                return url;
            }
        },
        props: ['pageSlug'],
        methods: {
            loadAsyncPageData: function() {
                var pageSlug = this.pageSlug;
                // fix url if needed...
                if(pageSlug == null) {
                    pageSlug = "";
                }

                this.$http.get(this.$store.state.backendURL + '/api/page/' + pageSlug).then(
                res => {
                    this.error = 0;
                    this.page = res.body;
                    this.renderBreadcrumb();
                }, res => {
                    this.error = res.status;
                        this.renderBreadcrumb();
                });
            },
            renderBreadcrumb: function() {
                var pageSlug = this.PageSlug;
                if(pageSlug != null) {
                    pageSlug = pageSlug.split("/");
                } else {
                    pageSlug = [];
                }

                var htmlList = [];
                htmlList.push("<li><a><i class='fa fa-home'></i></a></li>");

                // set breadcrumb
                this.breadcrumb = htmlList.join("");
            },
            redirectToPage: function(homepage) {
                this.$router.push({
                    name: "page",
                    params: {
                        spaceKey: this.spaceKey,
                        pageSlug: homepage
                    }
                })
            },
            deletePage() {
                this.$http.delete(
                    this.$store.state.backendURL + '/api/wiki/' + this.spaceKey + '/' + this.pageSlug
                ).then(
                    () => {
                        this.$toast.open({
                            type: 'is-success',
                            message: 'Page deleted!'
                        });

                        this.loadAsyncPageData();
                    }, res => {
                        this.error = res.status;
                    });
            }
        },
        mounted: function() {
            this.loadAsyncPageData();
        },
        watch: {
            '$route' (to, from) {
                if(to !== from) this.loadAsyncPageData();
            }
        }
    };
</script>