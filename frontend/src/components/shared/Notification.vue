<template>
    <section class="section"v-if="notification && notification.type">
        <div class="container">
            <b-message :title="notification.type === 'success' ? 'Success' : 'Error'"
                       :type="'is-' + notification.type"
                       :closable="false">
                {{ notification.message }}
            </b-message>
        </div>
    </section>
</template>

<script>
  export default {
    computed: {
      notification() {
        return this.$store.state.notification;
      }
    }
  };
</script>