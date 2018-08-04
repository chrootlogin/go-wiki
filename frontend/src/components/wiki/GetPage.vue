<template >
    <div v-if="error === 0">
        <article>
            <section class="hero is-primary">
                <div class="hero-body">
                    <div class="container">
                        <h1 v-if="!edit" class="title">
                            {{ page.title }}
                        </h1>
                        <div v-if="edit" class="field">
                            <div class="control">
                                <input class="input is-large is-primary" type="text" v-model="pageForm.title" placeholder="Large input">
                            </div>
                        </div>
                    </div>
                </div>
            </section>
            <div v-if="!edit" class="notification">
                <div class="container has-text-right">
                    <button class="button is-success">
                        <span>Create page</span>
                        <span class="icon is-small">
                            <i class="fa fa-plus"></i>
                        </span>
                    </button>
                    <button v-on:click="enableEdit" class="button is-primary">
                        <span>Edit page</span>
                        <span class="icon is-small">
                            <i class="fa fa-edit"></i>
                        </span>
                    </button>
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
            <div class="container" v-if="edit">
                <b-tabs type="is-boxed" v-on:change="editorTabChanged" v-model="activeEditorTab">
                    <b-tab-item label="Editor">
                        <vue-editor :disabled="disableWysiwygEditor" ref="editor" v-model="pageForm.content"></vue-editor>
                    </b-tab-item>
                    <b-tab-item label="Source">
                        <b-field label="Source code">
                            <b-input type="textarea" v-model="pageForm.content"></b-input>
                        </b-field>
                    </b-tab-item>
                    <b-tab-item label="Preview">
                        <div class="content" v-html="pageForm.content"></div>
                    </b-tab-item>
                </b-tabs>
                <p class="has-text-right">
                    <button v-on:click="cancelEdit" class="button is-danger">
                        <span>Cancel</span>
                        <span class="icon is-small">
                            <i class="fa fa-times"></i>
                        </span>
                    </button>
                    <button v-on:click="saveEdit" class="button is-primary" v-bind:class="{'is-loading': loading}">
                        <span>Save page</span>
                        <span class="icon is-small">
                            <i class="fa fa-save"></i>
                        </span>
                    </button>
                </p>
            </div>
        </article>
    </div>
</template>

<script>
    import { VueEditor } from 'vue2-editor'

    export default {
        components: {
            VueEditor
        },
        data() {
            return {
                edit: false,
                loading: false,
                disableWysiwygEditor: false,
                error: 0,
                page: {
                    title: "",
                    content: ""
                },
                activeEditorTab: 0,
                pageForm: {}
            }
        },
        props: ['pageSlug'],
        methods: {
            loadAsyncPageData: function() {
                // fix url if needed...
                if(this.pageSlug == null) {
                    this.pageSlug = "";
                }
                this.$http.get(this.$store.state.backendURL + '/api/page/' + this.pageSlug).then(
                res => {
                    this.error = 0;
                    console.log(res.body);
                    this.page = res.body;
                }, res => {
                    this.error = res.status;
                });
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
            enableEdit() {
                this.pageForm = JSON.parse(JSON.stringify(this.page));
                this.edit = true;
            },
            cancelEdit() {
                this.pageForm = {};
                this.edit = false;
            },
            saveEdit() {
                this.loading = true;

                this.$http.put(
                    this.$store.state.backendURL + '/api/wiki/' + this.spaceKey + '/' + this.pageSlug,
                    this.pageForm
                ).then(
                    () => {
                        this.loading = false;
                        this.edit = false;

                        this.$toast.open({
                            type: 'is-success',
                            message: 'Page updated!'
                        });

                        this.loadAsyncPageData();
                    }, res => {
                        this.loading = false;

                        this.error = res.status;
                    });
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
            },
            editorTabChanged(index) {
                if(index === 0) {
                    this.disableWysiwygEditor = false;
                } else {
                    this.disableWysiwygEditor = true;
                }
            }
        },
        mounted: function() {
            this.loadAsyncPageData();
        },
        watch: {
            "pageSlug": function(newVal, oldVal) {
                if(newVal !== oldVal) {
                    this.loadAsyncPageData();
                }
            }
        }
    };
</script>