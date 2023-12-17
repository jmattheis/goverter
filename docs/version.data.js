import util from "util";
import { exec as syncExec } from "child_process";
const exec = util.promisify(syncExec);

export default {
  async load() {
    try {
      return (await exec("git describe --tags --abbrev=0")).stdout
        .toString()
        .trim();
    } catch {
      return "unknown";
    }
  },
};
