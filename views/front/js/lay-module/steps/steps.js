("use strict");
layui.define("jquery", function (exports) {
  let $ = layui.jquery;
  /**
   * 缓存的实例对象
   */
  let thisModule = {};
  /**
   * 全局实例编号
   */
  let globalIndex = 1;
  const MAX_WIDTH = 800;

  /**
   * @namespace 用到的固定变量
   * @prop {string}  STEP_ACTION_KEY - 关键key的属性名称
   * @prop {string}  STEP_ACTION_BODY - 主体的class名称
   * @prop {string}  STEP_ACTION_PART - 主体的描述列表的class名称
   * @prop {string}  STEP_MARK_LABEL - 主体的标记是否展示label的属性
   * @prop {string}  STEP_MARK_STATUS - 主体的标记是否展示状态的属性
   * @prop {string}  STEP_MARK_DESC - 主体的标记是否展示描述的属性
   * @prop {string}  STEP_MARK_TOTAL - 主体的标记展示总步数的属性
   * @prop {string}  STEP_MARK_INDEX - 主体的标记指定步数的下标
   * @prop {string}  STEP_BAR_CLASS - 步骤条样式名称
   * @prop {string}  STEP_BAR_SELECTED - 被选中的步骤的样式class名称
   * @prop {string}  STEP_BAR_PART - 步骤条主体区域样式名称
   * @prop {string}  STEP_BAR_TEXT - 步骤条主体的文字区域样式名称
   * @prop {string}  STEP_BAR_EXTEND - 步骤条扩展属性(组)样式名称
   * @prop {string}  STEP_BAR_EXTEND_NAME - 步骤条扩展属性(名称)样式名称
   * @prop {string}  STEP_BAR_EXTEND_STATUS - 步骤条扩展属性(状态)样式名称
   * @prop {string}  STEP_BAR_EXTEND_CONTENT - 步骤条扩展属性(描述)样式名称
   * @prop {string}  STEP_BAR_EVENT_NEXT - 步骤条下一步按钮的样式名称
   * @prop {string}  STEP_BAR_EVENT_PREV - 步骤条上一步按钮的样式名称
   */
  const CONSTANT = {
    /**
     * @var 关键key的属性名称
     */
    STEP_ACTION_KEY: "layui-steps-key",
    /**
     * @var 主体的class名称
     */
    STEP_ACTION_BODY: "layui-steps-area",
    /**
     * @var 主体的描述列表的class名称
     * @desc 仅用于在初始化时进行扫描，将上面的参数压入处理的队列中
     */
    STEP_ACTION_PART: "layui-steps-part",
    /**
     * @var 主体的标记是否展示label的属性
     */
    STEP_MARK_LABEL: "lay-mark-label",
    /**
     * @var 主体的标记是否展示状态的属性
     */
    STEP_MARK_STATUS: "lay-mark-status",
    /**
     * @var 主体的标记是否展示描述的属性
     */
    STEP_MARK_DESC: "lay-mark-desc",
    /**
     * @var 主体的标记展示总步数的属性
     */
    STEP_MARK_TOTAL: "lay-mark-total",
    /**
     * @var 主体的标记指定步数的下标
     */
    STEP_MARK_INDEX: "lay-mark-index",
    /**
     * @var 步骤条样式名称
     */
    STEP_BAR_CLASS: "layui-steps-bar",
    /**
     * @var 被选中的步骤的样式class名称
     */
    STEP_BAR_SELECTED: "layui-steps-barSelected",
    /**
     * @var 步骤条主体区域样式名称
     */
    STEP_BAR_PART: "layui-steps-barPart",
    /**
     * @var 步骤条主体的文字区域样式名称
     */
    STEP_BAR_TEXT: "layui-steps-barText",
    /**
     * @var 步骤条扩展属性(组)样式名称
     */
    STEP_BAR_EXTEND: "layui-steps-barExtend",
    /**
     * @var 步骤条扩展属性(名称)样式名称
     */
    STEP_BAR_EXTEND_NAME: "layui-steps-extend-name",
    /**
     * @var 步骤条扩展属性(状态)样式名称
     */
    STEP_BAR_EXTEND_STATUS: "layui-steps-extend-status",
    /**
     * @var 步骤条扩展属性(描述)样式名称
     */
    STEP_BAR_EXTEND_CONTENT: "layui-steps-extend-content",
    /**
     * @var 步骤条下一步按钮的样式名称
     */
    STEP_BAR_EVENT_NEXT: "layui-steps-event-next",
    /**
     * @var 步骤条上一步按钮的样式名称
     */
    STEP_BAR_EVENT_PREV: "layui-steps-event-prev",
  };

  /**
   * @namespace 全局配置属性
   * @prop {string}  STEP_ACTION_BODY - 各步骤条是否启用自适应宽度
   *
   */
  const DEFAULT_CONFIG = {
    /**
     * @var 各步骤条是否启用自适应宽度
     * @desc
     *  此状态开启时，在调用render方法刷新步骤条是，回去检查最外层的
     * 容器当前的宽度和当前展示的步骤条个数，最终计算出一个步骤条的宽度，
     * 用这个结果来调整每个步骤条的宽度，从而达到占满整个容器的效果
     *  如果不开启，(这里统一开启) 那每个步骤条的最小宽度是80px
     */
    FIT: true,

    /**
     * 版本号
     */
    VERSION: "0.0.1",
  };

  /**
   * 状态枚举
   */
  const STATUS_ENUM = {
    SUCCESS: {
      id: "success",
      name: "已完成",
    },

    ERROR: {
      id: "error",
      name: "失败",
    },

    FINISH: {
      id: "finish",
      name: "已结束",
    },

    PROCESS: {
      id: "process",
      name: "进行中",
    },

    WAIT: {
      id: "wait",
      name: "就绪",
    },
  };

  /**
   * @constructor 步骤条实例
   */
  let worker = function (destination) {
    this.destination = destination;
    /**
     * 初始化步骤条的指针，当前指针为1
     */
    this.index = 1;
    /**
     * @method 获取当前容器宽度的计算值
     * @desc
     *    这里以 1 2 ...  7 8 9 ... 15 16 每个宽度80px
     * 计算得到，最小计算宽度是 800px 请保持容器的宽度不小于这个值
     * @returns 当前容器宽度
     */
    this._getWidth = function () {
      return this.destination.width() || MAX_WIDTH;
    };
    /**
     * @method 创建步骤条的容器
     * @desc
     *    一般的标签解析都要有步骤条容器，但是不排除后面扩展的js对象的声明方式，
     * 所以没有一个容器对象需要特殊创建
     */
    this._createContainer = function () {
      this.container = $(`<div class = "${CONSTANT.STEP_ACTION_BODY}"></div>`);
      this.destination.append(this.container);
    };
    /**
     * @method 初始化步骤条容器
     */
    this._initContainer = function () {
      let container = this.destination.find("." + CONSTANT.STEP_ACTION_BODY);
      if (container.length == 0) return this._createContainer();
      this.container = container.length == 1 ? container : $(container.get(0));
    };
    /**
     * @method 初始化配置项
     * @decs
     *
     *    source 当前步骤条的列表信息，页面根据这个信息进行渲染
     *    config 暂时规定下面几项
     *
     *  - LABEL  是否展示步骤名称
     *  - STATUS 是否展示当前状态
     *  - DESC  是否展示步骤描述
     *
     */
    this._initConfig = function () {
      this.source = [];
      this.config = {
        LABEL: false,
        STATUS: false,
        DESC: false,
      };
    };

    /**
     * @method 获取步骤条信息
     * @desc
     *    从dom的属性上面获取步骤条信息
     *
     * 1. 首先从步骤条容器上面获取信息(这个是指最外层容器里面的一个div，class名称为 {@linkplain CONSTANT.STEP_ACTION_BODY} ):
     *    - LABEL  是否展示步骤名称     对应属性名称:  {@linkplain CONSTANT.STEP_MARK_LABEL}
     *    - STATUS 是否展示当前状态     对应属性名称:  {@linkplain CONSTANT.STEP_MARK_STATUS}
     *    - DESC  是否展示步骤描述      对应属性名称:  {@linkplain CONSTANT.STEP_MARK_DESC}
     * 2. 遍历步骤条容器下面的 类 {@linkplain CONSTANT.STEP_ACTION_PART} 将属性放入 source 列表中
     *    - LABEL   步骤名称           对应属性名称: title
     *      description  步骤描述      是取的text() 值
     *      icon      代替layui图标    是取的这个dom下面的.layui-icon
     *      status    统一取就绪的状态
     *
     */
    this._scan = function () {
      // 获取 是否展示步骤名称 属性
      let labelFlag = this.container.attr(CONSTANT.STEP_MARK_LABEL);
      if (labelFlag != undefined && labelFlag == "true")
        this.config.LABEL = true;
      // 获取 是否展示当前状态 属性
      let statusFlag = this.container.attr(CONSTANT.STEP_MARK_STATUS);
      if (statusFlag != undefined && statusFlag == "true")
        this.config.STATUS = true;
      // 获取 是否展示步骤描述 属性
      let descFlag = this.container.attr(CONSTANT.STEP_MARK_DESC);
      if (descFlag != undefined && descFlag == "true") this.config.DESC = true;
      // 获取需要渲染的总步数
      let totalFlag = this.container.attr(CONSTANT.STEP_MARK_TOTAL);
      let self = this;
      if (!parseInt(totalFlag)) {
        // 遍历步骤条容器下面的 CONSTANT.STEP_ACTION_PART
        this.container.find("." + CONSTANT.STEP_ACTION_PART).each(function () {
          let $this = $(this);
          self.source.push({
            title: $this.attr("title") || "",
            description: $this.text() || "",
            icon: $this.find(".layui-icon").html(),
            status: STATUS_ENUM.WAIT,
          });
        });
      } else {
        // 首先渲染出相应数量的list元素值
        let total = parseInt(totalFlag);
        for (let i = 1; i <= total; i++) {
          self.source.push({
            title: "步骤" + i,
            description: "",
            icon: "",
            status: STATUS_ENUM.WAIT,
          });
        }

        let index = 1;
        // 遍历步骤条容器下面的 CONSTANT.STEP_ACTION_PART 补全信息
        this.container.find("." + CONSTANT.STEP_ACTION_PART).each(function () {
          let $this = $(this);
          let indexFlag = $this.attr(CONSTANT.STEP_MARK_INDEX);
          if (parseInt(indexFlag)) index = parseInt(indexFlag) - 1;
          self.source[index++] = {
            title: $this.attr("title") || "",
            description: $this.text() || "",
            icon: $this.find(".layui-icon").html(),
            status: STATUS_ENUM.WAIT,
          };
        });
      }
      // 确定当前的步骤条容器高度
      if (!DEFAULT_CONFIG.FIT) return;
      let totalHeight = 40;
      if (this.config.LABEL) totalHeight += 36;
      if (this.config.STATUS) totalHeight += 26;
      // 这个还是放在最后再统一计算吧
      // if (this.config.DESC) totalHeight += 150;
      this.totalHeight = totalHeight;
      this.container.css("height", totalHeight + "px");
    };

    /**
     * @method 生成步骤条
     * @returns 返回当前实例
     */
    this._buildBar = function () {
      /**
       * 根据当前的容器宽度和当前待渲染的列表个数来进行计算：
       *    按照一步100px
       *    如果能全部渲染，就调用全部渲染的方法进行渲染，
       * 如果不能就调用分页渲染
       */
      return this.source.length * 100 + 120 > this._getWidth()
        ? this._pageBar()
        : this._nomalBar();
    };

    /**
     * @method 全部渲染步骤条
     * @returns 返回当前实例
     */
    this._nomalBar = function () {
      // 调整步骤条的宽度
      if (DEFAULT_CONFIG.FIT) this._fitWidth(this.source.length);
      /**
       * 检查状态，如果之前使用分页渲染，这里需要清空了，重新渲染
       * 如果第一次没有进行渲染，这个是undefined那也是需要清空了重新渲染的
       */
      if (this.pageFlag == true || this.pageFlag == undefined) {
        this.container.empty();
        this.pageFlag = false;
        // 创建步骤条
        for (let i = 1; i <= this.source.length; i++) {
          this._appendBar(i);
        }
      } else {
        // 更新步骤条状态
        for (let i = 1; i <= this.source.length; i++) {
          this._updateBar(i);
        }
      }
      // 调整步骤条的高度
      if (DEFAULT_CONFIG.FIT) this._fitHeight();
      return this;
    };

    /**
     * @method 分页渲染步骤条
     * @returns 返回当前实例
     */
    this._pageBar = function () {
      // 不管是什么情况，分页渲染都是重新渲染，没得更新操作
      this.container.empty();
      this.pageFlag = true;
      // 计算当前的步骤条容器宽度最多支持多少个步骤同时渲染
      let total = Math.floor((this._getWidth() - 120) * 0.01);
      let length = this.source.length;
      // 如果可以同时渲染这么多个步骤就使用上面的全部渲染的方式。否则采用下面的分页渲染的方式
      if (length <= total) return this._nomalBar();
      // 当前的指针小于2 说明 进行中 这一步 出现在左边，使用左边渲染的方式
      if (this.index <= 2) return this._leftPageBar(total);
      // 当前的指针接近尾部 说明 进行中 这一步 出现在右边，使用右边渲染的方式
      if (this.index >= length - 2) return this._rightPageBar(total);
      // 使用中间渲染的方式
      return this._centerPageBar(total);
    };

    /**
     * @method 分页渲染步骤条(左)
     * @param {*} total  当前需要渲染的步骤个数
     * @desc 
     *    中间的格式必须是
     * 1 2 ... end
     *    要渲染4步，所以total不小于4
     * @returns 返回当前实例
     */
    this._leftPageBar = function (total) {
      // 分页特殊处理
      if(total < 4) total = 4; 
      // 调整步骤条的宽度
      if (DEFAULT_CONFIG.FIT) this._fitWidth(total);
      /**
       * 首先将左边的步骤添加进来 留两个位置
       */
      for (let i = 1; i <= total - 2; i++) {
        this._appendBar(i);
      }
      // 加一个 ... 未完成的状态
      this._appendBar("...", false);
      // 将最后一个放入进来
      this._appendBar(this.source.length);
      // 调整步骤条的高度
      if (DEFAULT_CONFIG.FIT) this._fitHeight();
      return this;
    };

    /**
     * @method 分页渲染步骤条(右)
     * @param {*} total 当前需要渲染的步骤个数
     * @desc 
     *    中间的格式必须是
     * 1 ... index index + 1 . end  1 ... 3 4 5   5-3+2
     *    要渲染end - index + 1 + 2 步 
     * @returns 返回当前实例
     */
    this._rightPageBar = function (total) {
      let conpare = this.source.length - this.index + 3;
      if(total < conpare) total = conpare;
      // 调整步骤条的宽度
      if (DEFAULT_CONFIG.FIT) this._fitWidth(total);
      // 渲染第一个步骤
      this._appendBar(1);
      // 加一个 ... 已完成的状态
      this._appendBar("...", true);
      // 将剩下的total-2个个数的步骤渲染出来
      for (let i = 1; i <= total - 2; i++) {
        this._appendBar(this.source.length - (total - 2) + i);
      }
      // 调整步骤条的高度
      if (DEFAULT_CONFIG.FIT) this._fitHeight();
      return this;
    };

    /**
     * @method 分页渲染步骤条(中)
     * @param {*} total 当前需要渲染的步骤个数
     * @desc 
     *    中间的格式必须是
     * 1 ... x x+1 x+2 ... end
     *    要渲染7步，所以total不小于7
     * @returns 返回当前实例
     */
    this._centerPageBar = function (total) {
      if(total < 7) total = 7;
      // 调整步骤条的宽度
      if (DEFAULT_CONFIG.FIT) this._fitWidth(total);
      // 渲染第一个步骤
      this._appendBar(1);
      // 加一个 ... 已完成的状态
      this._appendBar("...", true);
      // 计算中心位置
      let temp = Math.floor((total - 4) * 0.5);
      // 根据当前的指针位置再次调整这个渲染的中心位置
      if (this.index - temp < 2) temp = this._fitleft(temp, total);
      if (total - 4 + this.index - temp > this.source.length - 2)
        temp = this._fitright(temp, total);
      // 将中间 的 total -4 个步骤渲染出来
      for (let i = 1; i <= total - 4; i++) {
        this._appendBar(i + this.index - temp);
      }
      // 加一个 ... 未完成的状态
      this._appendBar("...", false);
      // 将最后一个放入进来
      this._appendBar(this.source.length);
      // 调整步骤条的高度
      if (DEFAULT_CONFIG.FIT) this._fitHeight();
      return this;
    };

    /**
     * @method 向左调整渲染中心，使进行中的步骤(当前指针位置)出现在渲染节点中
     * @param {*} temp 当前渲染中心位置
     * @param {*} total 当前需要渲染的步骤个数
     * @returns 调整后的渲染中心位置
     */
    this._fitleft = function (temp, total) {
      if (this.index - temp < 2) return this._fitleft(temp - 1, total);
      return temp;
    };

    /**
     * @method 向右调整渲染中心，使进行中的步骤(当前指针位置)出现在渲染节点中
     * @param {*} temp 当前渲染中心位置
     * @param {*} total 当前需要渲染的步骤个数
     * @returns 调整后的渲染中心位置
     */
    this._fitright = function (temp, total) {
      if (total - 4 + this.index - temp > this.source.length - 2)
        return this._fitright(temp + 1, total);
      return temp;
    };

    /**
     * @method 向渲染容器中加入一个执行的步骤
     * @param {*} pos 当前步骤在source中的位置 + 1
     * @param {*} flag 当前第一个参数是 ... 时，根据这个参数确定该节点是否已完成
     */
    this._appendBar = function (pos, flag) {
      // 判断是普通步骤还是 ...   true 就是普通节点
      let stepFlag = pos !== "...";
      // 获取当前步骤的状态
      let stepStatus = stepFlag
        ? this.source[pos - 1].status
        : flag
        ? STATUS_ENUM.SUCCESS
        : STATUS_ENUM.WAIT;
      // 获取当前的图标
      let iconHtml =
        stepFlag && this.source[pos - 1].icon ? this.source[pos - 1].icon : pos;
      this.container.append(
        $(`
        <div class = "${CONSTANT.STEP_BAR_CLASS} ${
          stepStatus === STATUS_ENUM.SUCCESS ||  stepStatus === STATUS_ENUM.PROCESS ? CONSTANT.STEP_BAR_SELECTED : ""
        }" pageindex = "${pos}" style = "width: ${this.barLength}px;">
          <div class = "${CONSTANT.STEP_BAR_PART}">
            <div class = "${CONSTANT.STEP_BAR_TEXT}">${iconHtml}</div>
          </div>
          ${
            stepFlag
              ? `
          <div class = "${CONSTANT.STEP_BAR_EXTEND}" style = "width: ${this.barLength}px; right: -${this.offsetRight}px;">
            ${
              this.config.LABEL
                ? `
            <div class = "${CONSTANT.STEP_BAR_EXTEND_NAME}">${
                    this.source[pos - 1].title
                      ? this.source[pos - 1].title
                      : `步骤${pos}`
                  }</div>
            `
                : ""
            }
            ${
              this.config.STATUS
                ? `
            <div class = "${CONSTANT.STEP_BAR_EXTEND_STATUS}">${stepStatus.name}</div>
            `
                : ""
            }
            ${
              this.config.DESC
                ? `
            <div class = "${CONSTANT.STEP_BAR_EXTEND_CONTENT}">${
                    this.source[pos - 1].description
                  }</div>
            `
                : ""
            }
          </div>
          `
              : ""
          }
        </div>
      `)
      );
    };

    /**
     * @method 更新当前步骤条dom的状态
     * @param {*} pos
     */
    this._updateBar = function (pos) {
      let $stepBar = this.destination.find('[pageindex="' + pos + '"]');
      let stepStatus = this.source[pos - 1].status;
      $stepBar
        .find("." + CONSTANT.STEP_BAR_EXTEND_STATUS)
        .text(stepStatus.name);
      $stepBar.css('width', this.barLength + 'px');
      $stepBar.find('.' + CONSTANT.STEP_BAR_EXTEND).css('width', this.barLength + 'px');
      // 注释掉,bug修复 https://gitee.com/giteetcj/layui_steps/issues/I7UNEE
      // $stepBar.css('right', '-' + this.offsetRight + 'px');  
      $stepBar.find('.' + CONSTANT.STEP_BAR_EXTEND).css('right', '-' + this.offsetRight + 'px');  
      if (stepStatus == STATUS_ENUM.SUCCESS || stepStatus == STATUS_ENUM.PROCESS) {
        if (!$stepBar.hasClass(CONSTANT.STEP_BAR_SELECTED))
          $stepBar.addClass(CONSTANT.STEP_BAR_SELECTED);
      } else {
        if ($stepBar.hasClass(CONSTANT.STEP_BAR_SELECTED))
          $stepBar.removeClass(CONSTANT.STEP_BAR_SELECTED);
      }
    };

    /**
     * @method 自适应步骤条的宽度
     * @param {*} total 总共渲染的步骤条数量
     */
    this._fitWidth = function (total) {
      if(!total) return;
      let width = Math.floor(this._getWidth() / (total + 1));
      let offset = Math.floor(width * 0.5) - 10;
      this.barLength = width;
      this.offsetRight = offset;
      // document.documentElement.style.setProperty(
      //   "--layui-steps-barLength",
      //   width + "px"
      // );
      // document.documentElement.style.setProperty(
      //   "--layui-steps-offsetRight",
      //   "-" + offset + "px"
      // );
    };
    /**
     * @method 自适应步骤条的高度
     */
    this._fitHeight = function(){
      if (this.config.DESC) {
        let totalHeight = this.totalHeight;
        let height = 0;
        this.container
          .find("." + CONSTANT.STEP_BAR_EXTEND_CONTENT)
          .each(function () {
            let contentHeight = $(this).height();
            if (contentHeight > height) height = contentHeight;
          });
        if (height)
          this.container.css("height", Number(totalHeight + height) + "px");
      }
    };
    /**
     * @method 添加点击事件
     */
    this._addListener = function () {
      let self = this;
      self.destination
        .off("click", "." + CONSTANT.STEP_BAR_EVENT_PREV)
        .on("click", "." + CONSTANT.STEP_BAR_EVENT_PREV, function () {
          self.prev();
        });
      self.destination
        .off("click", "." + CONSTANT.STEP_BAR_EVENT_NEXT)
        .on("click", "." + CONSTANT.STEP_BAR_EVENT_NEXT, function () {
          self.next();
        });
    };

    /**
     * @method 步骤条指针右移
     * @returns 当前实例
     */
    this.next = function () {
      if (this.index == this.source.length + 1) return;
      this.index++;
      return this.render();
    };

    /**
     * @method 步骤条指针左移
     * @returns 当前实例
     */
    this.prev = function () {
      if (this.index == 1) return;
      this.index--;
      return this.render();
    };

    /**
     * @method 步骤条指针回到第一步
     * @returns 当前实例
     */
    this.first = function () {
      this.index = 1;
      return this.render();
    };

    /**
     * @method 步骤条指针回到最后一步
     * @returns 当前实例
     */
    this.last = function () {
      this.index = this.source.length + 1;
      return this.render();
    };

    /**
     * @method 渲染步骤条
     * @returns 当前实例
     */
    this.render = function () {
      /**
       * 根据当前的指针，调整source里面的状态
       */
      for (let i = 1; i <= this.source.length; i++) {
        if (i == this.index) this.source[i - 1].status = STATUS_ENUM.PROCESS;
        if (i < this.index) this.source[i - 1].status = STATUS_ENUM.SUCCESS;
        if (i > this.index) this.source[i - 1].status = STATUS_ENUM.WAIT;
      }
      // 重新渲染步骤条
      return this._buildBar();
    };

    // 初始化容器
    this._initContainer();
    // 初始化配置
    this._initConfig();
    // 读取dom信息，初始化渲染数据
    this._scan();
    // 添加点击事件
    this._addListener();
    // 渲染步骤条
    this.render();
  };

  let hander = {
    /**
     * @method 通过实例的key获取实例对象
     * @param {*} key   实例的key，可以从dom的属性上面获取
     * @returns 返回实例对象
     */
    getInst: function (key) {
      let inst = thisModule[key];
      return inst instanceof worker ? inst : null;
    },

    /**
     * @method 创建并获取步骤条
     * @param {*} destination 指定一个jquery对象or一个id选择器
     * @returns 步骤条对象
     */
    build: function (destination) {
      /**
       * 如果是id选择器转化成jquery对象
       */
      if (Object.prototype.toString.call(destination) === "[object String]")
        destination = $(destination);
      /**
       * 获取dom上面的 {@linkplain CONSTANT.STEP_ACTION_KEY key属性值}
       *   - 如果取到了这个属性值，说明这个步骤条已经初始化过了，
       *        直接从 {@linkplain thisModule 缓存对象}中获取实例，调用render渲染方法返回实例
       *   - 如果没有获取到这个属性值，说明是第一次执行这个操作。首先将这个dom上面
       *        添加上key，然后创建一个worker对象，将这个对象和这个key放入缓存对象中，
       *        最后将worker对象返回
       */
      let key = destination.attr(CONSTANT.STEP_ACTION_KEY);
      //console.log(key);
      if (key !== undefined) {
        let inst = hander.getInst(key);
		//console.log(inst);
        return inst && inst instanceof worker ? inst.render() : null;
      }
      // 添加属性
      destination.attr(CONSTANT.STEP_ACTION_KEY, globalIndex);
      let inst = new worker(destination);
      // 放入缓存
      thisModule[globalIndex++] = inst;
      return inst;
    },

    /**
     * 渲染页面上全部符合的步骤条
     */
    buildAll: function(){
      let selectors = document.documentElement.querySelectorAll('.' + CONSTANT.STEP_ACTION_BODY);
      Array.prototype.forEach.call(selectors, function(selector){
        if(selector.parentElement) hander.build($(selector.parentElement));
      });
    },
  };

  /**
   * 立即渲染
   */
  layui.link(layui.cache.base + 'steps/steps.css');

  hander.buildAll();

  exports("steps", {
    getInst: hander.getInst,
    build: hander.build,
    buildAll: hander.buildAll,
    v: DEFAULT_CONFIG.VERSION,
  });
});
