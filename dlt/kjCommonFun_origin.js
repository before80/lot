var kjCommonFun = {
    curPage: 1,  //当前页
    nextPage: 1,  //要访问的页面
    ltType: '', //彩种类型
    ltName: '', //彩种类型名
    provinceId: 0,  //省份id，全国游戏为0 
    NoDataStr: '<div class="m-czNo" style="text-align: center;padding: 100px 0;"><img src="//static.sporttery.cn/res_1_0/jcw/images/ico_game_no.png" width="196" height="196"/></div>',
    postData: {
        "gameNo": "", "provinceId": "0", "pageSize": "30",
        "isVerify": commonV1Fun.comIsVerify, 'termLimits': 0
    }, //请求参数
    allPage: 0,  //总页数
    pageSize: 0,  //每页显示多少条
    startTerm: 0,  //开始期数
    endTerm: 0,  //结束期数
    termLimits: 0,  //近多少期
    callBackFun: '', //回调方法
    divHeight: 935,
    marryObj: {
        dlt: { startWith: 10008, endWith: 19080, addNum: 300000, gameId: 33800 },
        pls: { startWith: 10111, endWith: 19188, addNum: 320000, gameId: 28200 },
        plw: { startWith: '09359', endWith: 19188, addNum: 340000, gameId: 28300 },
        qxc: { startWith: '09154', endWith: 19081, addNum: 360000, gameId: 17100 },
        sfc: { startWith: '09110', endWith: 19086, addNum: 380000, gameId: 100 },
        rxj: { startWith: '09110', endWith: 19086, addNum: 400000, gameId: 29200 },
        bqc: { startWith: '09127', endWith: 19095, addNum: 420000, gameId: 31300 },
        jqc: { startWith: '09135', endWith: 19100, addNum: 440000, gameId: 31400 },
    },
    dfcMarryObj: {
        hlj61: { endWith: 19055, gameId: 14900 },
        fj367: { endWith: 19082, gameId: 11700 },
        fj317: { endWith: 19189, gameId: 11500 },
        fj317fj: { endWith: 19189, gameId: 10160 },
        fj225: { endWith: 19189, gameId: 11600 },
        zj205: { endWith: 19189, gameId: 9900 },
        zj61: { endWith: 19082, gameId: 10000 },
        hn41: { endWith: 19082, gameId: 49900 },
        js7: { endWith: 19109, gameId: 15900 },
    },
    /**
     * 非空判断
     *
     */
    kjCheckNull: function (nums) {
        if (nums == undefined || nums == null || nums == ''
            || nums == 'null' || nums == 'undefined') {
            return '';
        } else {
            return nums;
        }
    },
    /**
     * 非空判断
     *
     */
    kjCheckNullStr: function (nums) {
        if (nums == undefined || nums == null || nums == ''
            || nums == 'null' || nums == 'undefined') {
            return '--';
        } else {
            return nums;
        }
    },
    /**
     *  输出开奖列表页分页代码
     *  每页10个页码
        allNum  总页数
        curNum  当前页数
    */
    pageListType: function (allNum, curNum) {
        var bNum = 0;
        if (curNum % 10 === 0) {
            bNum = curNum - 9;
        } else {
            bNum = Math.floor(curNum / 10) * 10 + 1;
        }
        var pageStr = '<ul class="m-pager">'
            + '<li class="number u-pad10 ' + (curNum == 1 ? "no" : "") + '"  onclick="kjCommonFun.goNextPage(1)">首页</li>'
            + '<li class="number u-pad10 ' + (curNum == 1 ? "no" : "") + '"  onclick="kjCommonFun.goNextPage(' + ((curNum - 1) < 1 ? 1 : (curNum - 1)) + ')">&lt;&nbsp;上一页</li>';
        for (var i = bNum; i < (bNum + 10); i++) {
            if (i > allNum) {
                break;
            }
            pageStr += '<li class="number ' + (i == curNum ? "active" : "") + '" onclick="kjCommonFun.goNextPage(' + i.toString() + ')">' + i + '</li>';
        }
        pageStr += '<li class="number u-pad10 ' + (curNum == allNum ? "no" : "") + '" onclick="kjCommonFun.goNextPage(' + (curNum == allNum ? curNum : curNum + 1) + ')">下一页&nbsp;&gt;</li>'
        pageStr += '<li class="number u-pad10 ' + (curNum == allNum ? "no" : "") + '" onclick="kjCommonFun.goNextPage(' + allNum + ')">尾页</li></ul>';
        return pageStr;
    },
    /**
     * 翻页方法
     */
    goNextPage: function (p) {
        if (p == kjCommonFun.curPage) {
            return;
        }
        kjCommonFun.nextPage = p;
        kjCommonFun.curPage = p;
        //调用数据
        //  kjCommonFun.callDataCenter();
        kjCommonFun.iniParam();
    },
    /**
     * 调用接口数据
     * Fun 获取数据成功后调用的方法
     */

    getDrawApi: function (Fun) {
        var apiUrl = jsCommonDataV1.webApi ? jsCommonDataV1.webApi : '//webapi.sporttery.cn';
        $.ajax({
            url: apiUrl + '/gateway/lottery/getHistoryPageListV1.qry',
            type: 'get',
            dataType: "json",
            data: kjCommonFun.postData,
            success: function (data) {
                if (data.errorCode == 0) {
                    Fun(data.value);
                } else {
                    Fun([]);
                }
            }, error: function (data) {
                Fun([]);
            }
        })
    },
    urlGet: function () {
        var aQuery = window.location.href.split("?");
        var aGET = new Array();
        if (aQuery.length > 1) {
            var aBuf = aQuery[1].split("&");
            for (var i = 0, iLoop = aBuf.length; i < iLoop; i++) {
                var aTmp = aBuf[i].split("=");
                aGET[aTmp[0]] = aTmp[1];
            }
        }
        return aGET;
    },
    iniParam: function () {//初始化市获取页面参数

        // var getParam = kjCommonFun.urlGet();
        //获取当前页数
        if (kjCommonFun.curPage > 0) {
            kjCommonFun.postData.pageNo = kjCommonFun.curPage;
        }

        //页面展示条数
        if (kjCommonFun.pageSize == 10 || kjCommonFun.pageSize == 30 || kjCommonFun.pageSize == 50 || kjCommonFun.pageSize == 100) {
            kjCommonFun.postData.pageSize = kjCommonFun.pageSize;
        } else {
            kjCommonFun.postData.pageSize = 30;
        }
        //游戏类型
        kjCommonFun.postData.gameNo = kjCommonFun.ltType;
        //开始期数
        if (kjCommonFun.startTerm > 0) {
            kjCommonFun.postData.startTerm = kjCommonFun.startTerm;
        } else {
            if (kjCommonFun.postData.startTerm != undefined) {
                delete kjCommonFun.postData.startTerm;
            }
        }
        //结束期数
        if (kjCommonFun.endTerm > 0) {
            kjCommonFun.postData.endTerm = kjCommonFun.endTerm;
        } else {
            if (kjCommonFun.postData.endTerm != undefined) {
                delete kjCommonFun.postData.endTerm;
            }
        }
        if (kjCommonFun.termLimits > 0) {
            kjCommonFun.postData.termLimits = kjCommonFun.termLimits;
        } else {
            if (kjCommonFun.postData.termLimits != undefined) {
                delete kjCommonFun.postData.termLimits;
            }
        }
        kjCommonFun.callDataCenter();
    },

    /**
     * 调用接口数据
     */
    callDataCenter: function () {

        kjCommonFun.getDrawApi(kjCommonFun.callBackFun);
    },
    /**
     * 返回链接内容
     * gameNum 开奖期号
     * type 1 数字彩
     * pdfType 是否有pdf
     * isGetKjpdf 是否有开奖公告pdf
     * isGetXlpdf 是否有销量公告pdf
     * status 为true返回销量公告
     */
    linkGet: function (gameNum, type, pdfType, isGetKjpdf, isGetXlpdf, status) {

        if (type == 1) {
            return kjCommonFun.marryNum(gameNum, pdfType, isGetKjpdf, isGetXlpdf, status);
        } else {
            return kjCommonFun.dfcMarryNum(gameNum, pdfType, isGetKjpdf);
        }
    },
    /**
     * 返回链接内容
     * gameNum 开奖期号
     * type 1 数字彩
     * pdfType 是否有pdf
     * isGetKjpdf 是否有开奖公告pdf
     * isGetXlpdf 是否有销量公告pdf
     * status 为true返回销量公告
     */
     linkGetDescript: function (gameNum, pdfType, drawPdfUrl) {
        var type = kjCommonFun.ltName;
        var lotto = kjCommonFun.marryObj[type], num;
        if (pdfType == 2 || pdfType == 3) {//1为pdf
            num = lotto.addNum + parseInt(gameNum);
            var temUrl = kjCommonFun.urlGet();
            if (temUrl["url"] != undefined && temUrl["url"] != "") {
                if (temUrl["url"].indexOf('lottery.gov.cn') > -1 || temUrl["url"].indexOf('sporttery.cn') > -1) {
                    return '<a class="blue" target="_blank" href="' + temUrl["url"] + '/kj/kjxq.html?' + num + '">详情</a>';
                } else {
                    return '';
                }
            } else {
                return '';
            }
        } else if (gameNum > lotto.endWith) {
            if(drawPdfUrl == '' || drawPdfUrl == undefined){
                return " ";
            }else{
                return '<a class="blue" target="_blank" href="'+drawPdfUrl+'">详情</a>';
            }
        } else {
            return "";
        }
        
    },


    /**
     * 返回链接内容
     */
    linkGetOld: function (gameNum, type, status) {
        if (type == 1) {
            return kjCommonFun.marryNum(gameNum, status);
        } else {
            return kjCommonFun.dfcMarryNum(gameNum);
        }
    },
    /**
     * 数字彩玩法处理
     */
    marryNum: function (gameNum, pdfType, isGetKjpdf, isGetXlpdf, status) {
        var type = kjCommonFun.ltName, xl = status ? 'XL' : '';
        var lotto = kjCommonFun.marryObj[type], num;
        if (pdfType == 2 || pdfType == 3) {//1为pdf
            num = lotto.addNum + parseInt(gameNum);
            var temUrl = kjCommonFun.urlGet();
            if (temUrl["url"] != undefined && temUrl["url"] != "") {
                if (temUrl["url"].indexOf('lottery.gov.cn') > -1 || temUrl["url"].indexOf('sporttery.cn') > -1) {
                    return '<a class="blue" target="_blank" href="' + temUrl["url"] + '/kj/kjxq.html?' + num + '">详情</a>';
                } else {
                    return '';
                }
            } else {
                return '';
            }

        } else if (gameNum > lotto.endWith) {
            if (xl === 'XL' && isGetXlpdf == 0) {//未生成销量公告
                return "";
            } else if (xl === '' && isGetKjpdf == 0) {
                return "";
            } else {
                var ggid = lotto.gameId;
                if ($('.u-czList li').eq(1).hasClass("on")) {
                    ggid = '29200'
                }
                return '<a class="blue" target="_blank" href="https://pdf.sporttery.cn/' + ggid + '/' + gameNum + '/' + gameNum + '' + xl + '.pdf">详情</a>';
            }
        } else {
            return "";
        }
    },
    /**
     * 地方彩玩法处理
     */
    dfcMarryNum: function (gameNum, pdfType, isGetKjpdf) {
        var lotto = kjCommonFun.dfcMarryObj[kjCommonFun.ltName];
        if (gameNum >= lotto.endWith) {
            if (isGetKjpdf == 0) {
                return "";
            } else {
                return '<a class="blue" target="_blank" href="https://pdf.sporttery.cn/' + lotto.gameId + '/' + gameNum + '/' + gameNum + '.pdf">详情</a>';
            }
        } else {
            return "";
        }
    },
    /**
     * 页数及期数查询选择功能
     */

    pageAndTerm: function () {
        $(".u-zdy input").focus(function () {
            $(this).css({ "border-color": "#999" });
        });
        $(".u-zdy input").blur(function () {
            if ($(this).val() == '') {
                $(this).val(kjCommonFun.defaultValue);
                $(this).css({ "border-color": "#d8d8d8" });
            }
        });
        //显示页面条数选择
        $('.u-zjnum li').click(function () {
            if ($(this).attr("name") !== $("#showName").html()) {
                $('.u-zjnum li').removeClass('on');
                $(this).addClass('on');
                $("#showName").html($(this).attr("name"));
                kjCommonFun.termLimits = $(this).attr("showNum");
                kjCommonFun.pageSize = kjCommonFun.termLimits;

                kjCommonFun.startTerm = 0;
                kjCommonFun.endTerm = 0;
                $('.u-zjnum ul').hide();
                //  kjCommonFun.callDataCenter();
                kjCommonFun.curPage = 1;
                kjCommonFun.nextPage = 1;
                kjCommonFun.iniParam();
                $('#alertTxt').hide();
                $('#bterm').val('');
                $('#eterm').val('');
            }
        });
        $('.u-zjnum ').hover(function () {
            $('.u-zjnum ul').removeAttr("style");
        })
        $("#searchBrn").click(function () {
            var bNum = $('#bterm').val();
            var eNum = $('#eterm').val();
            if (bNum == "" && eNum == "") {
                kjCommonFun.pageSize = 0;
            } else {
                if (bNum == '') {
                    $('#alertTxt').text('请输入开始期号');
                    $('#alertTxt').show();
                    return;
                }
                if (eNum == '') {
                    $('#alertTxt').text('请输入结束期号');
                    $('#alertTxt').show();
                    return;
                }
                if (eNum <= bNum) {
                    $('#alertTxt').text('结束期号太小');
                    $('#alertTxt').show();
                    return;
                }
                if ((eNum.length == 5) && (bNum.length == 5)) {
                    if (parseInt(eNum.substr(0, 2) - parseInt(bNum.substr(0, 2))) > 5) {
                        $('#alertTxt').text('最多查询5年');
                        $('#alertTxt').show();
                        return;
                    }
                } else {
                    $('#alertTxt').text('输入正确的期号');
                    $('#alertTxt').show();
                    return;
                }
                kjCommonFun.pageSize = 30;
            }
            kjCommonFun.startTerm = bNum;
            kjCommonFun.endTerm = eNum;
            kjCommonFun.termLimits = 0;
            $("#showName").html('全部');
            $('.u-zjnum li').removeClass('on');
            $('.u-zjnum li').eq(0).addClass('on');
            kjCommonFun.curPage = 1;
            kjCommonFun.nextPage = 1;
            kjCommonFun.iniParam();
            $('.u-zdy ul').hide();
            $('#alertTxt').text('');
            $('#alertTxt').hide();
        });
        $('.u-zdy').hover(function () {
            $('.u-zdy ul').removeAttr("style");
        });
    },
    getUrlParam: function (name) {
        var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)"); //构造一个含有目标参数的正则表达式对象
        var r = window.location.search.substr(1).match(reg);  //匹配目标参数
        if (r != null) return unescape(r[2]); return null; //返回参数值
    },
    getIframeHeight: function (num) {
        var hh = parseInt(num);
        if (hh === "" || hh == null || isNaN(hh) || typeof (hh) == "undefined") {
            num = 911;
        }
        var iframeInfo = {
            message: num
        };
        window.parent.postMessage(iframeInfo, '*');
    },
    czOpenUrl: function (tid) {
        //传足页面历史开奖跳转
        //tid，1：胜负游戏 2：任选9场3：6场半全场4：4场进球 
        var pageUrl = {
            p1: 'sfc', p2: 'rxj', p3: 'bqc', p4: 'jqc'
        }
        var temUrl = kjCommonFun.urlGet();
        if (tid == "" || temUrl == undefined) {
            return false;
        }
        window.open(temUrl["url"] + '/kj/kjlb.html?' + pageUrl[tid]);
    },
    delOnmouseOverBg: function (div) {
        $(div).attr("style", "background:white");
    }
}
if (jsCommonDataV1 == undefined || jsCommonDataV1 == "") {
    var jsCommonDataV1 = JSON.parse(kjCommonFun.getUrlParam('cmData'));
}
/**
 * @description: 获取指纹ID及神策ID
 * @version: v1.0.0
 */
var checkParametersToWebapi ={
    'X-Sensors-ID':'',
    'X-Dev-FP':'',
    getVisitor:function(){
        var cookieCustomerVisitor = checkParametersToWebapi.getCookie('X-Dev-FP')
        if(cookieCustomerVisitor != ''){
            checkParametersToWebapi['X-Dev-FP'] = cookieCustomerVisitor
        }else{
            var localStorageVisitorID = localStorage.getItem('X-Dev-FP')
            if(localStorageVisitorID !=null){
                checkParametersToWebapi['X-Dev-FP'] = localStorageVisitorID
            }
        }
    },
    getDistinct:function(){
        var cookieSensor = checkParametersToWebapi.getCookie('X-Sensors-ID')
        checkParametersToWebapi['X-Sensors-ID'] = cookieSensor==''?'':cookieSensor
        var localStorageSensorsId = localStorage.getItem('X-Sensors-ID')
        if(checkParametersToWebapi['X-Sensors-ID'] == '' && localStorageSensorsId !=null){
            checkParametersToWebapi['X-Sensors-ID'] = localStorageSensorsId
        }
    },

    /**
     * 获取cookie
     * @param {String} cname 键名
     */
    getCookie: function getCookie(cname) {
        var name = cname + "=";
        var ca = document.cookie.split(';');

        for (var i = 0; i < ca.length; i++) {
            var c = ca[i];

            while (c.charAt(0) == ' ') {
                c = c.substring(1);
            }

            if (c.indexOf(name) == 0) {
                return c.substring(name.length, c.length);
            }
        }

        return "";
    },
}
// checkParametersToWebapi.getVisitor();
// checkParametersToWebapi.getDistinct();
