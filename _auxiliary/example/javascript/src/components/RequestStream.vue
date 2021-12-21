<template>
  <div>
    <button @click="callServerMethod">Request Stream</button>
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
      console.log("request - stream call...");
      if (this.socket) {
        const message = {ping: "say hello!"};
        this.socket
          .requestStream({
            data: message,
            metadata: {
              trace: "123",
              method: "/hello/hellosrv/saystream",
              service: "api.hello",
              version: "v1.0.0"
            }
          })
          .subscribe({
            onComplete: () => {
              console.log("requestStream done");
              this.received.push("requestStream done");
            },
            onError: error => {
              console.log("got error with requestStream");
              console.error(error);
            },
            onNext: value => {
              // console.log("got next value in requestStream..");
              this.received.push(value.data);
            },
            // Nothing happens until `request(n)` is called
            onSubscribe: sub => {
              console.log("subscribe request Stream!");
              sub.request(3);
              this.sent.push(message);
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