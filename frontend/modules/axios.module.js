const axios = require("axios");

const postApi = async (url, body, configPost) => {
  try {
    return await axios.post(url, body, configPost);
  } catch (e) {
    return e.message;
  }
}

const getApi = async (url, configGet) => {
  try {
    return await axios.get(url, configGet);
  } catch (e) {
    return e.message;
  }
}

const putApi = async (url, body, configPut) => {
  try {
    return await axios.put(url, body, configPut);
  } catch (e) {
    return e.message;
  }
}

const deleteApi = async (url, configDelete) => {
  try {
    return await axios.delete(url, configDelete);
  } catch (e) {
    return e.message;
  }
}

module.exports = {
  postApi,
  getApi,
  putApi,
  deleteApi,
}