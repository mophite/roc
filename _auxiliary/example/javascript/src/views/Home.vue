<template>
  <div class="home">
    <p>Welcome in the RSocket JavaScript tester</p>
    <p>RSocket is connected: {{ isConnected }}</p>

    <div class="responses">
      <RequestResponse :socket="socket"/>
      <hr/>
      <RequestStream :socket="socket"/>
      <hr/>
      <RequestChannel :socket="socket"/>
      <hr/>
      <FlowableTest/>
    </div>
  </div>
</template>

<script>
import {RSocketClient} from "rsocket-core";
import RSocketWebSocketClient from "rsocket-websocket-client";
// @ is an alias to /src
import RequestResponse from "@/components/RequestResponse.vue";
import RequestStream from "@/components/RequestStream.vue";
import RequestChannel from "@/components/RequestChannel.vue";
import FlowableTest from "@/components/FlowableTest.vue";
import {JsonSerializer} from "rsocket-core/build/RSocketSerialization";

export default {
  name: "home",
  components: {
    RequestResponse,
    RequestStream,
    FlowableTest,
    RequestChannel
  },
  data() {
    return {
      socket: null
    };
  },
  methods: {
    connect() {
      console.log("connecting with RSocket...");
      // const transport = new RSocketWebSocketClient({
      //   url: "ws://localhost:7777/hello"
      // });
      const client = new RSocketClient({
        // send/receive JSON objects instead of strings/buffers
        serializers: {
          data: JsonSerializer,
          metadata: JsonSerializer
        },
        setup: {
          // ms btw sending keepalive to server
          keepAlive: 60000,
          // ms timeout if no keepalive response
          lifetime: 180000,
          // format of `data`
          dataMimeType: "application/json",
          // format of `metadata`
          metadataMimeType: "application/json",
        },

        transport: new RSocketWebSocketClient(
            {
              debug: true,
              url: 'ws://localhost:7777/hello',
            },
        ),
      });
      client.connect().subscribe({
        onComplete: socket => {
          this.socket = socket;
        },
        onError: error => {
          console.log("got connection error");
          console.error(error);
        },
        onSubscribe: cancel => {
          /* call cancel() to abort */
        }
      });

      setTimeout(() => {
      }, 30000000);
    }
  },

  computed: {
    isConnected() {
      return !!this.socket;
    }
  },
  mounted() {
    this.connect();
  }
};
</script>

<style scoped>
div.responses {
  width: 700px;
  margin: auto;
}
</style>
