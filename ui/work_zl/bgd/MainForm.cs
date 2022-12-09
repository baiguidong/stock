/*
 * 由SharpDevelop创建。
 * 用户： Administrator
 * 日期: 2018/11/4
 * 时间: 19:34
 * 
 * 要改变这种模板请点击 工具|选项|代码编写|编辑标准头文件
 */
using System;
using System.Collections.Generic;
using System.Drawing;
using System.Windows.Forms;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Diagnostics;
using System.IO;
using System.Threading;

namespace bgd
{
	/// <summary>
	/// Description of MainForm.
	/// </summary>
	public partial class MainForm : Form
	{
		public MainForm()
		{
			//
			// The InitializeComponent() call is required for Windows Forms designer support.
			//
			InitializeComponent();
			
			numericUpDown1.Value=0;
//			numericUpDown3.Value=0;
//			numericUpDown4.Value=0;
			comboBox1.Text="60";
			comboBox1.Items.Add("30");
			comboBox1.Items.Add("60");
			comboBox1.Items.Add("120");
			comboBox1.Items.Add("240");
			
			dateTimePicker1.Format=DateTimePickerFormat.Custom;
			dateTimePicker1.CustomFormat="yyyy-MM-dd";

		}
		
		//输出原图2
		void Button2Click(object sender, EventArgs e)
		{
			string jb=" -jb="+comboBox1.Text.Trim();
			string py=" -py="+numericUpDown1.Value;
			string day=" -day="+dateTimePicker1.Text.Trim();
			string gd=" -gd="+numericUpDown4.Value;
			string kd=" -kd="+numericUpDown3.Value;
			string xt=" -xt="+numericUpDown2.Value;
			
			string para=day+jb+py+gd+kd+xt+" -ty=10";
			System.Diagnostics.Process exep = System.Diagnostics.Process.Start(@"svr_1.exe",
                                                                               @para);
			exep.WaitForExit();
		}
		void MainFormLoad(object sender, EventArgs e)
		{
	
		}
		//输出原图
		void Button1Click(object sender, EventArgs e)
		{
			string jb=" -jb="+comboBox1.Text.Trim();
			string py=" -py="+numericUpDown1.Value;
			string gd=" -gd="+numericUpDown4.Value;
			string kd=" -kd="+numericUpDown3.Value;
			string xt=" -xt="+numericUpDown2.Value;
			string day=" -day="+dateTimePicker1.Text.Trim();
			
			string para=day+jb+py+gd+kd+xt+" -ty=1";

			System.Diagnostics.Process exep = System.Diagnostics.Process.Start(@"svr_1.exe",
                                                                               @para);
			exep.WaitForExit();
		}
		void Button3Click(object sender, EventArgs e)
		{
//			string jb=" -jb="+comboBox1.Text.Trim();
//			string py=" -py="+numericUpDown1.Value;
//			string day=" -day="+dateTimePicker1.Text.Trim();
//			string gd=" -gd="+numericUpDown4.Value;
//			string kd=" -kd="+numericUpDown3.Value;
//			string xt=" -xt="+numericUpDown2.Value;
//			
//			string para=day+jb+py+gd+kd+xt+" -ty=4";
//			System.Diagnostics.Process exep = System.Diagnostics.Process.Start(@"svr_1.exe",
//                                                                               @para);
//			exep.WaitForExit();

// 输出新图2
			string jb=" -jb="+comboBox1.Text.Trim();
			string py=" -py="+numericUpDown1.Value;
			string py2=" -py2="+numericUpDown5.Value;
			string x=" -x="+numericUpDown6.Value;
			string gd=" -gd="+numericUpDown4.Value;
			string kd=" -kd="+numericUpDown3.Value;
			string xt=" -xt="+numericUpDown2.Value;
			string day=" -day="+dateTimePicker1.Text.Trim();
			
			string para=day+jb+py+py2+x+gd+kd+xt+" -ty=2";;

			System.Diagnostics.Process exep = System.Diagnostics.Process.Start(@"svr_2.exe",
                                                                               @para);
			exep.WaitForExit();
			
		}
		void Button4Click(object sender, EventArgs e)
		{
			
			string jb=" -jb="+comboBox1.Text.Trim();
			string py=" -py="+numericUpDown1.Value;
			string day=" -day="+dateTimePicker1.Text.Trim();
			string gd=" -gd="+numericUpDown4.Value;
			string kd=" -kd="+numericUpDown3.Value;
			string xt=" -xt="+numericUpDown2.Value;
			string gm=" -gm="+numericUpDown6.Value;
			
			string para=day+jb+py+gd+kd+xt+gm+" -ty=4";
			System.Diagnostics.Process exep = System.Diagnostics.Process.Start(@"svr_1.exe",
                                                                               @para);
			exep.WaitForExit();

//			string jb=" -jb="+comboBox1.Text.Trim();
//			string py=" -py="+numericUpDown1.Value;
//			string day=" -day="+dateTimePicker1.Text.Trim();
//			string gd=" -gd="+numericUpDown4.Value;
//			string kd=" -kd="+numericUpDown3.Value;
//			string xt=" -xt="+numericUpDown2.Value;
//			
//			string para=day+jb+py+gd+kd+xt+" -ty=3";
//			System.Diagnostics.Process exep = System.Diagnostics.Process.Start(@"svr_1.exe",
//                                                                               @para);
//			exep.WaitForExit();
		}
		void Button5Click(object sender, EventArgs e)
		{
			string jb=" -jb="+comboBox1.Text.Trim();
			string py=" -py="+numericUpDown1.Value;
			string py2=" -py2="+numericUpDown5.Value;
			string x=" -x="+numericUpDown6.Value;
			string gd=" -gd="+numericUpDown4.Value;
			string kd=" -kd="+numericUpDown3.Value;
			string xt=" -xt="+numericUpDown2.Value;
			string day=" -day="+dateTimePicker1.Text.Trim();
			string ks=" -ks="+numericUpDown7.Text.Trim();
			string gs=" -gs="+numericUpDown8.Text.Trim();
			string xx=" -xx="+numericUpDown9.Text.Trim();
			
			string para=day+jb+py+py2+x+gd+kd+xt+ks+gs+xx;

			System.Diagnostics.Process exep = System.Diagnostics.Process.Start(@"svr_2.exe",
                                                                               @para);
			exep.WaitForExit();
		}
		void NumericUpDown7ValueChanged(object sender, EventArgs e)
		{
	
		}
		void NumericUpDown6ValueChanged(object sender, EventArgs e)
		{
	
		}
		void NumericUpDown9ValueChanged(object sender, EventArgs e)
		{
	
		}
		void Button6Click(object sender, EventArgs e)
		{
			string jb=" -jb="+comboBox1.Text.Trim();
			string py=" -py="+numericUpDown1.Value;
			string day=" -day="+dateTimePicker1.Text.Trim();
			string gd=" -gd="+numericUpDown4.Value;
			string kd=" -kd="+numericUpDown3.Value;
			string xt=" -xt="+numericUpDown2.Value;
			string gm=" -gm="+numericUpDown6.Value;
			
			string para=day+jb+py+gd+kd+xt+gm+" -ty=6";
			System.Diagnostics.Process exep = System.Diagnostics.Process.Start(@"svr_1.exe",
                                                                               @para);
			exep.WaitForExit();
		}
	}
}
