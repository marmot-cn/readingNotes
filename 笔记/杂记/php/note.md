[首页](/) > [编程目录](/code/) > [php代码目录](/code/php/)> 代码注释规范

#代码注释规范

---

####通常，代码应该被详细地注释。这不仅仅有助于给缺乏经验的程序员描述代码的流程和意图，而且有助于给你提供丰富的内容以让你在几个月后再看自己的代码时仍能很好的理解。 注释没有强制规定的格式，但是我们建议以下的形式。

文档块(DocBlock) 式的注释要写在类和方法的声明前，这样它们就能被集成开发环境(IDE)捕获：


    <?php
    /**
     * Super Class
     *
     * @package Package Name
     * @subpackage Subpackage
     * @category Category
     * @author Author Name
     * @link http://example.com
     */
    class Super_class {

    }

    /**
     * Encodes string for use in XML
     *
     * @access public
     * @param string
     * @return string
     */
    function xml_encode($str)

使用行注释时，在大的注释块和代码间留一个空行。

    <?php
    // break up the string by newlines
    $parts = explode("\n", $str);

    // A longer comment that needs to give greater detail on what is
    // occurring and why can use multiple single-line comments.  Try to
    // keep the width reasonable, around 70 characters is the easiest to
    // read.  Don't hesitate to link to permanent external resources
    // that may provide greater detail:
    //
    // http://example.com/information_about_something/in_particular/

    $parts = $this->foo($parts);

所有的注释标示

    <?php
    /**
     * @name 为某个变量指定别名
     * @abstract 记录一个抽象类，类变量或方法
     * @access 属性的访问、使用权限. @access private 表明这个属性是被保护的。
     * @author 文档作者的名字和邮箱地址
     * @category 组织packages
     * @copyright 指明版权信息
     * @deprecated 文档中被废除的方法
     * @const 指明常量
     * @deprecate 指明不推荐或者是废弃的信息
     * @example 文档的外部保存的示例文件的位置
     * @exception 文档中方法抛出的异常，也可参照 @throws.
     * @global 文档中的全局变量及有关的方法和函数
     * @ignore 忽略文档中指定的关键字
     * @internal 开发团队内部信息
     * @link 类似于license 但还可以通过link找到文档中的更多个详细的信息
     * @exclude 指明当前的注释将不进行分析，不出现在文挡中
     * @final 指明这是一个最终的类、方法、属性，禁止派生、修改。
     * @global 指明在此函数中引用的全局变量
     * @include 指明包含的文件的信息
     * @link 定义在线连接
     * @module 定义归属的模块信息
     * @modulegroup 定义归属的模块组
     * @package 定义归属的包的信息
     * @param 定义函数或者方法的参数信息
     * @return 定义函数或者方法的返回信息
     * @see 文件关联的任何元素（全局变量，包括，页面，类，函数，定义，方法，变量）。
     * @since 指明该api函数或者方法是从哪个版本开始引入的
     * @static 指明变量、类、函数是静态的。
     * @throws 指明此函数可能抛出的错误异常,极其发生的情况
     * @todo 指明应该改进或没有实现的地方
     * @var 定义说明变量/属性。
     * @version 定义版本信息
     */

注释实例：

    <?php
      /**
      * start page for webaccess
      *
      * PHP version 5
      *
      * @category  PHP
      * @package   PSI_Web
      * @author    Michael Cramer <BigMichi1@users.sourceforge.net>
      * @copyright 2009 phpSysInfo
      * @license   http://opensource.org/licenses/gpl-2.0.php GNU General Public License
      * @version   SVN: $Id: class.Webpage.inc.php 412 2010-12-29 09:45:53Z Jacky672 $
      * @link      http://phpsysinfo.sourceforge.net
      */
      /**
      * generate the dynamic webpage
      *
      * @category  PHP
      * @package   PSI_Web
      * @author    Michael Cramer <BigMichi1@users.sourceforge.net>
      * @copyright 2009 phpSysInfo
      * @license   http://opensource.org/licenses/gpl-2.0.php GNU General Public License
      * @version   Release: 3.0
      * @link      http://phpsysinfo.sourceforge.net
      */
     class Webpage extends Output implements PSI_Interface_Output
     {
         /**
          * configured language
          *
          * @var String
          */
         private $_language;

         /**
          * configured template
          *
          * @var String
          */
         private $_template;

         /**
          * all available templates
          *
          * @var Array
          */
         private $_templates = array();

         /**
          * all available languages
          *
          * @var Array
          */
         private $_languages = array();

         /**
          * check for all extensions that are needed, initialize needed vars and read config.php
          */
         public function __construct()
         {
             parent::__construct();
             $this->_getTemplateList();
             $this->_getLanguageList();
         }

         /**
          * checking config.php setting for template, if not supportet set phpsysinfo.css as default
          * checking config.php setting for language, if not supported set en as default
          *
          * @return void
          */
         private function _checkTemplateLanguage()
         {
             $this->_template = trim(PSI_DEFAULT_TEMPLATE);
             if (!file_exists(APP_ROOT.'/templates/'.$this->_template.".css")) {
                 $this->_template = 'phpsysinfo';
             }

             $this->_language = trim(PSI_DEFAULT_LANG);
             if (!file_exists(APP_ROOT.'/language/'.$this->_language.".xml")) {
                 $this->_language = 'en';
             }
         }

         /**
          * get all available tamplates and store them in internal array
          *
          * @return void
          */
         private function _getTemplateList()
         {
             $dirlist = CommonFunctions::gdc(APP_ROOT.'/templates/');
             sort($dirlist);
             foreach ($dirlist as $file) {
                 $tpl_ext = substr($file, strlen($file) - 4);
                 $tpl_name = substr($file, 0, strlen($file) - 4);
                 if ($tpl_ext === ".css") {
                     array_push($this->_templates, $tpl_name);
                 }
             }
         }

         /**
          * get all available translations and store them in internal array
          *
          * @return void
          */
         private function _getLanguageList()
         {
             $dirlist = CommonFunctions::gdc(APP_ROOT.'/language/');
             sort($dirlist);
             foreach ($dirlist as $file) {
                 $lang_ext = substr($file, strlen($file) - 4);
                 $lang_name = substr($file, 0, strlen($file) - 4);
                 if ($lang_ext == ".xml") {
                     array_push($this->_languages, $lang_name);
                 }
             }
         }

         /**
          * render the page
          *
          * @return void
          */
         public function run()
         {
             $this->_checkTemplateLanguage();

             $tpl = new Template("/templates/html/index_dynamic.html");

             $tpl->set("template", $this->_template);
             $tpl->set("templates", $this->_templates);
             $tpl->set("language", $this->_language);
             $tpl->set("languages", $this->_languages);

             echo $tpl->fetch();
         }
      }

更详细的注释规范可以查看[phpDocumentor](http://manual.phpdoc.org/HTMLSmartyConverter/HandS/phpDocumentor/tutorial_phpDocumentor.howto.pkg.html#basics.docblock)



#私有注释规范

---

此处罗列我们暂时使用到的一些标签规范的细节备注

**@version**

		v.大版本.小版本.时间
		
		大版本:从1开始
		小版本:从0开始
		时间:当前时间:年月日
		
		v1.0.20150527		

**@author**
		
		写上自己的git账号名称

**@example**
	
		为了使用自文档性,请在注释里面备写上示例.如果有必要请在连接上 .md 格式文件
		
**redmine**

		对所有注释请明确标记redmine任务编号. 以 #xxx 格式.



		