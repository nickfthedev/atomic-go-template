import axios from "axios";
import { readFile, writeFile } from "fs";

// Make a request for a user with a given ID
axios
  .get("http://127.0.0.1:4040/api/tunnels")
  .then(function (response) {
    // handle success
    const publicURL = response.data.tunnels[0].public_url;
    console.log("--> ONLY FOR DEV MODE! NGROK URL: ", publicURL);

    // Open and read .env file
    readFile(".env", "utf8", (err, data) => {
      if (err) {
        console.error(err);
        return;
      }
      let searchString = "SHOPIFY_APP_URL=";
      let re = new RegExp("^.*" + searchString + ".*$", "gm");
      let formatted = data.replace(re, "SHOPIFY_APP_URL=" + publicURL);

      writeFile(".env", formatted, "utf8", function (err) {
        if (err) return console.log(err);
      });
      console.log("--> .env file has been updated! ");
    });
  })
  .catch(function (error) {
    // handle error
    console.log(error);
  })
  .finally(function () {
    // always executed
  });
