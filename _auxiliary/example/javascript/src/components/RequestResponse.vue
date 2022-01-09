<template>
  <div>
    <button @click="callServerMethod">Request Response</button>
    <SentReceivedViewer
        :sent="sent"
        :received="received"
    />
  </div>
</template>

<script>
import SentReceivedViewer from "@/components/SentReceivedViewer.vue";

export default {
  props: {
    socket: {
      type: Object,
      default: null
    }
  },
  components: {
    SentReceivedViewer
  },
  data() {
    return {
      sent: [],
      received: []
    };
  },
  methods: {
    callServerMethod() {
      console.log("request - response call...");
      if (this.socket) {
        this.socket
            .requestResponse({
              data: {ping: "ping"},
              metadata: {
                trace: "123",
                method: "/hello/hello/sayapi",
                service: "api.hello",
                version: "v1.0.0",
              }
            })
            .subscribe({
              onComplete: data => {
                console.log("got response with requestResponse");
                this.received.push(data.data);
              },
              onError: error => {
                console.log("got error with requestResponse");
                console.error(error);
              },
              onSubscribe: cancel => {
                this.sent.push({ping: "ping"});
                /* call cancel() to stop onComplete/onError */
              }
            });
      } else {
        console.log("not connected...");
      }
    }
  }
};
</script>

<style lang="scss" scoped>
</style>