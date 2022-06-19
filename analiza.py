import matplotlib.pyplot as plt
#Helper file to create some of the plots.


# line 1 points
x1 = [2015,2016,2017,2018,2019,2020,2021,2022]
y1 = [0.1550,0.2095,0.2552,0.2400,0.2571,0.2642,0.2548,0.2811]
# plotting the line 1 points 
plt.plot(x1, y1, label = "libc6")

# line 1 points
x2 = [2015,2016,2017,2018,2019,2020,2021,2022]
y2 = [0.1253,0.1689,0.2202,0.1970,0.2260,0.2218,0.2243,0.1906]
# plotting the line 1 points 
plt.plot(x2, y2, label = "libgcc1")

# line 1 points
x3 = [2015,2016,2017,2018,2019,2020,2021,2022]
y3 = [0.0789,0.0907,0.0772,0.0571,0.0583,0.0500,0.0457,0.0387]
# plotting the line 1 points 
plt.plot(x3, y3, label = "multiarch-support")

# line 1 points
x4 = [2015,2016,2017,2018,2019,2020,2021,2022]
y4 = [0.0148,0.,0.0071,0.0060,0.0054,0.0055,0.0058,0.0063]
# plotting the line 1 points 
plt.plot(x4, y4, label = "zlib1g")

# line 1 points
x5 = [2015,2016,2017,2018,2019,2020,2021,2022]
y5 = [0.0050,0.0074,0.0066,0.0062,0.0079,0.0065,0.0080,0.0070]
# plotting the line 1 points 
plt.plot(x5, y5, label = "dpkg")

x6 = [2015,2016,2017,2018,2019,2020,2021,2022]
y6 = [0.0139,0.0075,0.0074,0.0152,0.0149,0.0121,0.0107,0.0087]
# plotting the line 1 points 
plt.plot(x6, y6, label = "gcc")

x7 = [2015,2016,2017,2018,2019,2020,2021,2022]
y7 = [0.0092,0.0069,0.0083,0.0093,0.0082,0.0081,0.0082,0.0072]
# plotting the line 1 points 
plt.plot(x7, y7, label = "perl")

x8 = [2018,2019,2020,2021,2022]
y8 = [0.0193,0.0040,0.0165,0.0217,0.0132]
# plotting the line 1 points 
plt.plot(x7, y7, label = "libjs-jquery")


  
# # line 2 points
# x2 = [1,2,3]
# y2 = [4,1,3]
# # plotting the line 2 points 
# plt.plot(x2, y2, label = "line 2")
  
# naming the x axis
plt.xlabel('Year')
# naming the y axis
plt.ylabel('PageRank')
# giving a title to my graph
#plt.title('Two lines on same graph!')
  
# show a legend on the plot
plt.legend()
  
# function to show the plot
plt.show()