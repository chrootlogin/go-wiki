<template>
    <section class="section" v-if="error === 0">
        <div class="container">
            <div class="notification">
                <nav class="breadcrump is-hidden-mobile" aria-label="breadcrumbs">
                    <ul v-html="breadcrumb"></ul>
                </nav>
            </div>
            <div class="columns">
                <div class="column is-one-fifth">
                    <aside class="menu">
                        <p class="menu-label">
                            Page Administration
                        </p>
                        <ul class="menu-list">
                            <li>
                                <router-link :to="{ name: 'edit', query: { path: url } }">
                                    <span>Edit page</span>
                                    <span class="icon is-small">
                                        <i class="fa fa-edit"></i>
                                    </span>
                                </router-link>
                            </li>
                            <li>
                                <a v-on:click="createPage">
                                    <span>Create page</span>
                                    <span class="icon is-small">
                                        <i class="fa fa-plus"></i>
                                    </span>
                                </a>
                            </li>
                            <li>
                                <a v-on:click="deletePage">
                                    <span>Delete page</span>
                                    <span class="icon is-small">
                                        <i class="fa fa-times"></i>
                                    </span>
                                </a>
                            </li>
                        </ul>
                    </aside>
                </div>
                <div class="column">
                    <article class="content" v-html="page.content"></article>
                </div>
            </div>
        </div>
    </section>
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
            createPage() {
                this.$dialog.prompt({
                    message: `URL of the new page?`,
                    inputAttrs: {
                        placeholder: 'e.g. my_page/',
                        maxlength: 255
                    },
                    onConfirm: (value) => {
                        this.$router.push({
                            name: "edit",
                            query: {
                                path: value
                            }
                        })
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
            console.log(this.url)
            this.loadAsyncPageData();
        },
        watch: {
            '$route' (to, from) {
                if(to !== from) this.loadAsyncPageData();
            }
        }
    };
</script>