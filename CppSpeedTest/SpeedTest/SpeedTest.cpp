

#include <iostream>


#include <cstdlib>
#include <string>
#include <vector>
#include <ctime>
#include <numeric>
#include <cmath>
#include <sstream>

#include <thread>
#include <chrono>
#include <mutex>
#include <thread>
#include <bit>

using namespace std;


std::vector<unsigned char> writeArray(int size)
{
    //unsigned char buffer[1024];
    std::cout << "Starting to write Array" << endl;

    std::vector<unsigned char> byteArray;

    for (int i = 0; i <= size; i++)
    {
        byteArray.push_back(rand() % 2);
    }

    std::cout << "Done writing array" << endl;
    return byteArray;
}

std::vector<unsigned char> transfer(vector<unsigned char> input)
{
    std::vector<unsigned char> byteArrayTransfer;

    byteArrayTransfer = input;
    for (unsigned i = 1; i < input.size(); i++)
    {
        byteArrayTransfer.push_back(input[i]);
    }

    return byteArrayTransfer;
}


int main()
{
    static int size = 1024; //Setting the size of the arrays

    std::cout << "Size of the Array is set to: " << size << endl;

    static std::vector<unsigned char> bA1;
    static std::vector<unsigned char> bA2;
    static std::vector<unsigned char> bA3;

    auto startTime = chrono::steady_clock::now();

    std::thread t1([] {bA1 = writeArray(size);});
    std::thread t2([] {bA2 = writeArray(size);});
    std::thread t3([] {bA3 = writeArray(size);});



   
    t1.join();
    t2.join();
    t3.join();

    auto endTime = chrono::steady_clock::now();

    std::cout << "Done writing arrays. Duration: " << chrono::duration_cast<chrono::milliseconds>(endTime - startTime).count() << "ms" <<endl;
    
    auto startTime1 = chrono::steady_clock::now();

    std::thread t4([] {bA2 = transfer(bA1); });
    std::thread t5([] {bA3 = transfer(bA2); });
    std::thread t6([] {bA1 = transfer(bA3); });

    t4.join();
    t5.join();
    t6.join();

    auto endTime1 = chrono::steady_clock::now();

    std::cout << "\nDone transferring arrays." << endl << "Duration: " << chrono::duration_cast<chrono::milliseconds>(endTime1 - startTime1).count() 
              << "ms" << "\nSpeed: " << size * (chrono::duration_cast<chrono::milliseconds>(endTime1 - startTime1).count()) / 60 
              << " mBit/s" << endl;
}

