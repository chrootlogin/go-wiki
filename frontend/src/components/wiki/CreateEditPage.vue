<template >
    <div v-if="error === 0">
        <section class="hero is-primary">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title">
                        Editor
                    </h1>
                </div>
            </div>
        </section>
        <div class="container">
            <b-tabs type="is-boxed" v-on:change="editorTabChanged" v-model="activeEditorTab">
                <b-tab-item label="Editor">
                    <markdown-editor v-model="page.content" ref="markdownEditor"></markdown-editor>
                </b-tab-item>
                <b-tab-item label="Preview">
                    <div class="content" v-html="preview"></div>
                </b-tab-item>
                <b-tab-item label="Settings">
                    <b-field label="Path">
                        <b-input v-model="path" disabled></b-input>
                    </b-field>
                </b-tab-item>
            </b-tabs>
            <p class="has-text-right">
                <button v-on:click="cancel" class="button is-danger">
                    <span>Cancel</span>
                    <span class="icon is-small">
                        <i class="fa fa-times"></i>
                    </span>
                </button>
                <button v-on:click="save" class="button is-primary" v-bind:class="{'is-loading': loading}">
                    <span>Save page</span>
                    <span class="icon is-small">
                        <i class="fa fa-save"></i>
                    </span>
                </button>
            </p>
        </div>
    </div>
</template>

<script>
    import markdownEditor from 'vue-simplemde/src/markdown-editor.vue'

    export default {
        components: {
            markdownEditor
        },
        data() {
            return {
                loading: false,
                error: 0,
                page: {
                    title: "",
                    content: ""
                },
                activeEditorTab: 0,
                preview: ""
            }
        },
        computed: {
            path() {
                let path = this.$route.query.path;

                // fix url if needed...
                if(path == null) {
                    path = "";
                }

                return path;
            }
        },
        methods: {
            loadAsyncPageData: function() {
                this.$http.get(this.$store.state.backendURL + '/api/page/' + this.path + '?format=no-render').then(
                res => {
                    this.error = 0;
                    this.page = res.body;
                }, res => {
                    this.error = res.status;
                });
            },
            loadAsyncPreviewData: function() {
                this.$http.post(
                    this.$store.state.backendURL + '/api/preview',
                    {"content": this.page.content}
                ).then(
                    res => {
                        this.error = 0;
                        this.preview = res.body.content;
                    }, res => {
                        this.error = res.status;
                    }
                );
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
            cancel() {
                this.redirectToPage(this.path);
            },
            save() {
                this.loading = true;

                this.$http.put(
                    this.$store.state.backendURL + '/api/page/' + this.path,
                    this.page
                ).then(
                    () => {
                        this.loading = false;

                        this.$toast.open({
                            type: 'is-success',
                            message: 'Page updated!'
                        });

                        this.redirectToPage(this.path);
                    }, res => {
                        this.loading = false;

                        this.error = res.status;
                    });
            },
            editorTabChanged(index) {
                if(index === 1) {
                    this.loadAsyncPreviewData();
                }
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

<style>
    @import '~simplemde/dist/simplemde.min.css';
</style>