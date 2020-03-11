
# include imports needed (some are commented out which will be needed for further implementation)
import os
import sys
import ember
#import numpy as np
import pandas as pd
#import pefile as pe
from numba import jit
import lightgbm as lgb
#import matplotlib.pyplot as plt
#from reports import Report, HTMLTable

# Remove this later
import warnings
warnings.filterwarnings("ignore")

# speed up computation using jit (cannot run in nopython mode and try-catch requires python 3.7, potentially find new solution)
@jit(parallel = True)
def buildPathList(path,pathList,extension, subFolders):
    for entry in os.scandir(path):
        if entry.is_file() and entry.path.endswith(extension):
            pathList.append(entry.path)
        elif entry.is_dir() and subFolders:   # if its a directory, then repeat process as a nested function
            pathList = findFilesInFolder(entry.path, pathList, extension, subFolders) 

# Function to recursively search through a directory for a folder of a certain type
# ** may need to change this implementation to work better with large directories **
# Obtained solution from: https://stackoverflow.com/questions/3964681/find-all-files-in-a-directory-with-extension-txt-in-python
def findFilesInFolder(path, pathList, extension, subFolders = True):
    """  Recursive function to find all files of an extension type in a folder (and optionally in all subfolders too)

    path:        Base directory to find files
    pathList:    A list that stores all paths
    extension:   File extension to find
    subFolders:  Bool.  If True, find files in all subfolders under path. If False, only searches files in the specified folder
    """
    try:   # Trapping a OSError:  File permissions problem I believe
        buildPathList(path, pathList, extension, subFolders)
    except OSError:
        print('Cannot access ' + path +'. Probably a permissions error')

    return pathList

# populate table data for html report
@jit(parallel = True)
def getTableData(model, pathListEXE, pathListDLL):
    predData = []

    for file in pathListEXE:
        test_data = open(file, "rb").read()

        predData.append({'File Name': os.path.basename(file), 'Prediction': "Malware" if int(round(ember.predict_sample(model, test_data))) == 1 else "Benign"})
        #print(os.path.basename(file))
  
    for file in pathListDLL:
        test_data = open(file, "rb").read()

        predData.append({'File Name': os.path.basename(file), 'Prediction': "Malware" if int(round(ember.predict_sample(model, test_data))) == 1 else "Benign"})
        #print(os.path.basename(file))

    # create dataframe from prediction struct
    predDF = pd.DataFrame(predData)
    return predDF

def main():  
    try:
        dir_name = sys.argv[1]
    except IndexError:
        print("ERROR: No directory given")
        return
    
    # load in Ember prediction model
    lgbm_model = lgb.Booster(model_file="ember_model_2018.txt")

    # create .exe path list and populate it
    exePathList = []
    exePathList = findFilesInFolder(dir_name, exePathList, ".exe", True)

    # create .dll path list and populate it
    dllPathList = []
    dllPathList = findFilesInFolder(dir_name, dllPathList, ".dll", True)

    # get model predictions from path lists
    pred_df = getTableData(lgbm_model,exePathList,dllPathList)
    
    # convert the dataframe into HTML
    html_table = pred_df.to_html(justify = 'left')
    
    # re-format HTML to include cell background color based on prediction
    lines = []
    for line in html_table.splitlines():
        lnSplit = line.split('<td')
        if 'Benign' in line:
            line = lnSplit[0] + '<td bgcolor="lightgreen"' + lnSplit[1]
        elif 'Malware' in line:
            line = lnSplit[0] + '<td bgcolor="#ff6666"' + lnSplit[1]
        lines.append(line)
    html_table = "".join(lines)

    # write HTML to aggregate report file
    Html_file= open("aggregate.html","w")
    Html_file.write(html_table)
    Html_file.close()

    return

if __name__ == "__main__":
    main()